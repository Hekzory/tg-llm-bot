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
	config             *config.ServiceConfig
	messageToDbQueue   chan tgbotapi.Message
	messageToUserQueue chan models.Message
	actionInChat       chan models.Message
}

func NewTelegramHandler(service *service.UserService, logger *logging.Logger, config *config.ServiceConfig) *TelegramHandler {
	return &TelegramHandler{
		service:            service,
		logger:             logger,
		config:             config,
		messageToDbQueue:   make(chan tgbotapi.Message, 100),
		messageToUserQueue: make(chan models.Message, 100),
		actionInChat:       make(chan models.Message, 100),
	}
}

func (h *TelegramHandler) StartServer(port int) {
	h.logger.Info("Server starts!")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	telegramBot := h.tgBotInitialisation()

	go h.processDatabase(ctx)
	go h.pollTelegram(ctx, telegramBot)
	go h.processChatAction(ctx)

	select {}
}

func (h *TelegramHandler) processDatabase(ctx context.Context) {
	ticker := time.NewTicker(h.config.Config.DBPollingInterval)
	defer ticker.Stop()
	for {
		select {

		case <-ctx.Done():
			h.logger.Info("Stopping telegram polling due to context cancellation")
			return

		case <-ticker.C:
			messages, err := h.service.GetReadyMessages(ctx)
			if err != nil {
				h.logger.Error("Error fetching new messages: %s", err)
				continue
			}
			for _, message := range messages {
				h.messageToUserQueue <- message
			}
			h.logger.Debug("Polling db done")

		case message := <-h.messageToDbQueue:

			if err := h.addNewUser(ctx, message); err != nil {
				h.logger.Error("Error adding new user: %s", err)
			}

			id, err := h.service.GetUserIdByTgId(ctx, int(message.Chat.ID))
			if err != nil {
				h.logger.Error("Error getting user by tg_id: %s", err)
			}

			if err := h.service.AddMessage(ctx, id, message.Text); err != nil {
				h.logger.Error("Error sending message to db: %s", err)
			}
			h.logger.Info("Message sent")

		}
	}

}

func (h *TelegramHandler) pollTelegram(ctx context.Context, telegramBot *telegram.Bot) {

	updates := h.GetUpdatesChannel(telegramBot)
	for {
		select {

		case <-ctx.Done():
			h.logger.Info("Stopping telegram polling due to context cancellation")
			return

		case message := <-h.messageToUserQueue:
			tg_id, err := h.service.GetTgIdByUserId(ctx, message.UserID)
			if err != nil {
				h.logger.Error("Error getting tg id: %s", err)
			}

			msg := tgbotapi.NewMessage(int64(tg_id), message.Answer)

			_, err = telegramBot.Bot.Send(msg)
			if err != nil {
				h.logger.Error("Error sending meassege to user: %s", err)
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

		case message := <-h.actionInChat:
			h.sendActionToUserID(ctx, telegramBot, message.UserID, "typing")
		}
	}
}

func (h *TelegramHandler) processChatAction(ctx context.Context) {
	ticker := time.NewTicker(h.config.TGSConfig.ActionPollingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			meassages, err := h.service.GetProcessingMessages(ctx)
			if err != nil {
				h.logger.Error("Error getting processing msg: %s", err)
			}
			if meassages != nil {
				for _, meassage := range meassages {
					h.actionInChat <- meassage
				}
			}
		}
	}
}

func (h *TelegramHandler) tgBotInitialisation() *telegram.Bot {
	bot, err := tgbotapi.NewBotAPI(h.config.TGSConfig.TelegramToken)
	if err != nil {
		h.logger.Fatal("Error creating new bot api: %s", err)
	}
	telegramBot := telegram.NewBot(bot)
	h.logger.Info("Authorized an account %s", telegramBot.Bot.Self.UserName)
	return telegramBot
}

func (h *TelegramHandler) addNewUser(ctx context.Context, message tgbotapi.Message) error {
	exist, err := h.service.UserExists(ctx, int(message.Chat.ID))
	if err != nil {
		h.logger.Error("Error checking is user exist: %s", err)
		return err
	}

	if !exist {
		err := h.service.AddUser(ctx, int(message.Chat.ID), message.From.FirstName, message.From.UserName)
		if err != nil {
			h.logger.Error("Error adding user to db %s", err)
			return err
		}
	}
	return nil
}

func (h *TelegramHandler) GetUpdatesChannel(telegramBot *telegram.Bot) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates, err := telegramBot.Bot.GetUpdatesChan(u)
	if err != nil {
		h.logger.Error("Error getting updates: %s", err)
	}
	return updates
}

func (h *TelegramHandler) sendActionToUserID(ctx context.Context, telegramBot *telegram.Bot, id int, action_type string) {
	tg_id, err := h.service.GetTgIdByUserId(ctx, id)
	if err != nil {
		h.logger.Error("Error getting tg id: %s", err)
	}
	action := tgbotapi.NewChatAction(int64(tg_id), "typing")
	telegramBot.Bot.Send(action)
}
