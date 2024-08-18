package handler

import (
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/config"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/service"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/telegram"
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramHandler struct {
	service            *service.UserService
	logger             *logging.Logger
	config             *config.Config
	messageToDbQueue   chan tgbotapi.Message
	messageToUserQueue chan models.Message
}

func NewTelegramHandler(service *service.UserService, logger *logging.Logger, config *config.Config) *TelegramHandler {
	return &TelegramHandler{
		service:            service,
		logger:             logger,
		config:             config,
		messageToDbQueue:   make(chan tgbotapi.Message, 100),
		messageToUserQueue: make(chan models.Message, 100),
	}
}

func (h *TelegramHandler) StartServer(port int) {
	h.logger.Info("Server starts!")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		h.logger.Fatal("Error creating new bot api: %s", err)

	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, h.logger)

	go h.processDatabase(ctx)
	go h.pollTelegram(ctx, telegramBot)

	select {}
}

func (h *TelegramHandler) processDatabase(ctx context.Context) {
	ticker := time.NewTicker(h.config.DBpollingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Stopping telegram polling due to context cancellation")
			return
		case <-ticker.C:
			messages, err := h.service.GetReadyMessages(ctx)
			if err != nil {
				h.logger.Fatal("Error fetching new messages: %s", err)
				continue
			}
			for _, message := range messages {
				h.messageToUserQueue <- message
			}
			h.logger.Debug("Polling db done")
		case message := <-h.messageToDbQueue:
			exist, err := h.service.IsUserExist(ctx, int(message.Chat.ID))
			if err != nil {
				h.logger.Fatal("Error checking is user exist: %s", err)
			}
			if !exist {
				h.service.AddUser(ctx, int(message.Chat.ID))
			}

			id, err := h.service.GetUserIdByTgId(ctx, int(message.Chat.ID))
			if err != nil {
				h.logger.Fatal("Error getting user by tg_id: %s", err)
			}

			err = h.service.AddMessage(
				ctx,
				id,
				message.Text)
			if err != nil {
				h.logger.Fatal("Error sending message to db: %s", err)
			}
			h.logger.Info("Message sent")

		}
	}

}

func (h *TelegramHandler) pollTelegram(ctx context.Context, telegramBot *telegram.Bot) {
	h.logger.Info("Authorized an account %s", telegramBot.Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telegramBot.Bot.GetUpdatesChan(u)
	if err != nil {
		h.logger.Fatal("Error getting updates: %s", err)
	}
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Stopping telegram polling due to context cancellation")
			return
		case message := <-h.messageToUserQueue:
			h.logger.Info("%d", message.UserID)
			tg_id, err := h.service.GetTgIdByUserId(ctx, message.UserID)
			if err != nil {
				h.logger.Fatal("Error getting tg id: %s", err)
			}
			h.logger.Info("%D", tg_id)
			msg := tgbotapi.NewMessage(int64(tg_id), message.Answer)
			_, err = telegramBot.Bot.Send(msg)
			if err != nil {
				h.logger.Fatal("Error sending meassege to user: %s", err)
			}
			if err := h.service.UpdateMessageStatus(ctx, message.ID, "answered"); err != nil {
				h.logger.Error("Error updating message status: %s", err)
				continue
			}

		case update := <-updates:
			if update.Message == nil {
				continue
			}
			h.messageToDbQueue <- *update.Message

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//telegramBot.Bot.Send(msg)
		}

	}
}
