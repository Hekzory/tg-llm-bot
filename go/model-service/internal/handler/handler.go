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
	service   *service.UserService
	logger    *logging.Logger
	cfg       *config.Config
	userQueue chan models.User
}

func NewModelHandler(service *service.UserService, logger *logging.Logger, config *config.Config) *ModelHandler {
	return &ModelHandler{
		service:   service,
		logger:    logger,
		cfg:       config,
		userQueue: make(chan models.User, 100),
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
	go h.processNewUsers(ctx)

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
		var response map[string]interface{}
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
		users, err := h.service.GetAllUsers(ctx)
		if err != nil {
			h.logger.Error("Error fetching new users: %s", err)
			continue
		}

		for _, user := range users {
			h.userQueue <- user // Добавляем пользователя в канал
		}

		h.logger.Debug("Polling db done")

		time.Sleep(h.cfg.DBPollingInterval)
	}
}

func (h *ModelHandler) processNewUsers(ctx context.Context) {
	for {
		select {
		case user := <-h.userQueue: // Читаем пользователя из канала
			h.logger.Debug(fmt.Sprintf("Got user to process: %v", user))
			go h.processUser(ctx, user)
		case <-ctx.Done():
			return
		}
	}
}

func (h *ModelHandler) processUser(ctx context.Context, user models.User) {
	url := h.cfg.ModelApiUrl + "api/generate"
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  h.cfg.ModelName,
		"prompt": "hi, my name is " + user.Username,
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

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		h.logger.Error("Error unmarshalling response body: %s, Response: %s", err, response)
		return
	}

	welcomeMessage, ok := response["response"].(string)
	if !ok {
		h.logger.Error("Message not found in response")
		h.logger.Error(fmt.Sprintf("%v", response))
		return
	} else {
		h.logger.Info(welcomeMessage)
	}
}
