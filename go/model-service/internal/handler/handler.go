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
	"time"
)

type ModelHandler struct {
	service              *service.ModelService
	logger               *logging.Logger
	cfg                  *config.Config
	messageQuestionQueue chan models.Message
	messageResultQueue   chan models.Message
}

func NewModelHandler(service *service.ModelService, logger *logging.Logger, config *config.Config) *ModelHandler {
	return &ModelHandler{
		service:              service,
		logger:               logger,
		cfg:                  config,
		messageQuestionQueue: make(chan models.Message, 100),
		messageResultQueue:   make(chan models.Message, 100),
	}
}

func (h *ModelHandler) StartServer() {
	h.logger.Info("Server starts!")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pull the model before starting the server
	err := h.pullModel(ctx)
	if err != nil {
		h.logger.Fatal("Failed to pull model: %s", err)
		return
	}

	go h.pollDatabase(ctx)
	go h.processNewMessages(ctx)

	select {} // Block forever
}

func (h *ModelHandler) pullModel(ctx context.Context) error {
	url := h.cfg.ModelApiUrl + "api/pull"
	requestBody, err := json.Marshal(map[string]string{
		"model": h.cfg.ModelName,
	})
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request to Ollama API: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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

		status, ok := response["status"].(string)
		if ok && status == "success" {
			h.logger.Info("Model pulled successfully")
			return nil
		} else {
			//h.logger.Info(response["status"].(string))
		}
	}

	return fmt.Errorf("model pull failed, no success status received")
}

func (h *ModelHandler) pollDatabase(ctx context.Context) {
	for {
		select {
		case message := <-h.messageResultQueue:
			err := h.service.UpdateMessage(ctx, message)
			if err != nil {
				h.logger.Error("Error updating message: %s", err)
			}
		default:
			messages, err := h.service.GetNewMessages(ctx)
			if err != nil {
				h.logger.Error("Error fetching new messages: %s", err)
				continue
			}

			for _, message := range messages {
				err := h.service.UpdateMessageStatus(ctx, message.ID, "processing")
				if err != nil {
					h.logger.Error("Error updating message status: %s", err)
					continue
				}
				h.messageQuestionQueue <- message
			}

			h.logger.Debug("Polling db done")
			time.Sleep(h.cfg.DBPollingInterval)
		}
	}
}

func (h *ModelHandler) processNewMessages(ctx context.Context) {
	for {
		select {
		case message := <-h.messageQuestionQueue:
			h.logger.Debug(fmt.Sprintf("Got message to process: %v", message))
			go h.processNewMessage(ctx, message)
		case <-ctx.Done():
			return
		}
	}
}

func (h *ModelHandler) processNewMessage(ctx context.Context, message models.Message) {
	url := h.cfg.ModelApiUrl + "api/generate"
	requestBody, err := json.Marshal(map[string]any{
		"model":  h.cfg.ModelName,
		"prompt": message.Question,
		"stream": false,
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

	client := &http.Client{}
	resp, err := client.Do(req)
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
		h.logger.Error(fmt.Sprintf("Message not found in response %v", response))
		return
	}
}
