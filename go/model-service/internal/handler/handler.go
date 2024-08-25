package handler

import (
	"Hekzory/tg-llm-bot/go/model-service/internal/config"
	"Hekzory/tg-llm-bot/go/model-service/internal/service"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ModelHandler struct {
	service              *service.ModelService
	logger               *logging.Logger
	cfg                  *config.ServiceConfig
	messageQuestionQueue chan models.Message
	messageResultQueue   chan models.Message
	httpClient           *http.Client
}

func NewModelHandler(service *service.ModelService, logger *logging.Logger, config *config.ServiceConfig) *ModelHandler {
	return &ModelHandler{
		service:              service,
		logger:               logger,
		cfg:                  config,
		messageQuestionQueue: make(chan models.Message, 100),
		messageResultQueue:   make(chan models.Message, 100),
		httpClient:           &http.Client{Timeout: config.ModelConfig.ModelAnswerTimeout},
	}
}

// ModelRequest represents the request body for model-related HTTP requests.
type ModelRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (h *ModelHandler) StartServer() {
	h.logger.Info("Server starts!")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Run pullModel concurrently
	go func() {
		defer wg.Done()
		if err := h.pullModel(ctx); err != nil {
			h.logger.Fatal("Failed to pull model: %s", err)
		}
	}()

	// Run resetStuckMessages concurrently
	go func() {
		defer wg.Done()
		if err := h.resetStuckMessages(ctx); err != nil {
			h.logger.Fatal("Failed to reset stuck messages: %s", err)
		}
	}()

	// Wait for both init tasks to complete before running the main loop
	wg.Wait()

	wg.Add(2)
	go func() { defer wg.Done(); h.pollDatabase(ctx) }()
	go func() { defer wg.Done(); h.processNewMessages(ctx) }()

	// Implement graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	h.logger.Info("Received signal to exit. Shutting down...")

	cancel() // Cancel the context to stop the goroutines

	// Wait for the goroutines to finish
	wg.Wait()

	h.logger.Info("Server stopped.")
}

func (h *ModelHandler) resetStuckMessages(ctx context.Context) error {
	messages, err := h.service.GetStuckMessages(ctx, h.cfg.ModelConfig.ModelAnswerTimeout)
	if err != nil {
		return fmt.Errorf("error fetching stuck messages: %w", err)
	}

	for _, message := range messages {
		if err := h.service.UpdateMessageStatus(ctx, message.ID, "new"); err != nil {
			h.logger.Error("Error resetting message status: %s", err)
		}
	}

	h.logger.Info("Stuck messages reset successfully")
	return nil
}

func (h *ModelHandler) pullModel(ctx context.Context) error {
	url := h.cfg.ModelConfig.ModelApiUrl + "api/pull"
	requestBody, err := json.Marshal(map[string]string{
		"model": h.cfg.ModelConfig.ModelName,
	})
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request to Ollama API: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to pull model: %s", body)
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var response map[string]any
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding response body: %w", err)
		}

		if status, ok := response["status"].(string); ok && status == "success" {
			h.logger.Info("Model pulled successfully")
			return nil
		}
	}

	return fmt.Errorf("model pull failed, no success status received")
}

func (h *ModelHandler) pollDatabase(ctx context.Context) {
	ticker := time.NewTicker(h.cfg.Config.DBPollingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Stopping database polling due to context cancellation")
			return
		case <-ticker.C:
			messages, err := h.service.GetNewMessages(ctx)
			if err != nil {
				h.logger.Error("Error fetching new messages: %s", err)
				continue
			}
			for _, message := range messages {
				if err := h.service.UpdateMessageStatus(ctx, message.ID, "processing"); err != nil {
					h.logger.Error("Error updating message status: %s", err)
					continue
				}
				h.messageQuestionQueue <- message
			}
			h.logger.Debug("Polling db done")
		case message := <-h.messageResultQueue:
			err := h.service.UpdateMessage(ctx, message)
			if err != nil {
				h.logger.Error("Error updating message: %s", err)
			}
		}
	}
}

func (h *ModelHandler) processNewMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Stopping message processing due to context cancellation")
			return
		case message := <-h.messageQuestionQueue:
			h.logger.Debug("Got message to process: %v", message)
			go h.processNewMessage(ctx, message)
		}
	}
}

func (h *ModelHandler) processNewMessage(ctx context.Context, message models.Message) {
	url := h.cfg.ModelConfig.ModelApiUrl + "api/generate"

	requestBody, err := json.Marshal(ModelRequest{
		Model:  h.cfg.ModelConfig.ModelName,
		Prompt: message.Question,
		Stream: false,
	})

	if err != nil {
		h.logger.Error("Error marshalling request body: %s", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		h.logger.Error("Error creating request to Ollama API: %s", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Error("Error making request to Ollama API: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("Error reading response body: %s", err)
		return
	}

	var response map[string]any
	err = json.Unmarshal(body, &response)
	if err != nil {
		h.logger.Error("Error unmarshalling response body: %s, Response: %s", err, response)
		return
	}

	answer, ok := response["response"].(string)
	if ok {
		message.Answer = answer
		h.messageResultQueue <- message
	} else {
		h.logger.Error("Message not found in response %v", response)
		return
	}
}
