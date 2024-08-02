package handler

import (
	"Hekzory/tg-llm-bot/go/telegram-service/internal/service"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

type ModelHandler struct {
	service *service.ModelService
	logger  *logging.Logger
}

func NewModelHandler(service *service.ModelService, logger *logging.Logger) *ModelHandler {
	return &ModelHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ModelHandler) StartServer(port int) {
	h.logger.Info("Server starts!")

	h.logger.Info("Server stops!")
}
