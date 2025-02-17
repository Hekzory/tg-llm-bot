package handler

import (
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/config"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/service"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/telegram"
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

			userId, err := h.service.GetUserIdByTgId(ctx, int(message.Chat.ID))
			if err != nil {
				h.logger.Error("Error getting user id: %s", err)
			}

			if err := h.startNewConversation(ctx, userId); err != nil {
				h.logger.Error("Error starting new conversation: %s", err)
			}

			convId, err := h.service.GetConvIdByUserId(ctx, userId)
			if err != nil {
				h.logger.Error("Error getting user id: %s", err)
			}

			if err, _ := h.service.AddMessage(ctx, message.Text, convId, message.MessageID); err != nil {
				h.logger.Error("Error sending message to db: %s", err)
			}
			h.logger.Info("Message sent")

		}
	}

}

func (h *TelegramHandler) pollTelegram(ctx context.Context, telegramBot *telegram.Bot) {

	updates := h.getUpdatesChannel(telegramBot)
	for {
		select {

		case <-ctx.Done():
			h.logger.Info("Stopping telegram polling due to context cancellation")
			return

		case message := <-h.messageToUserQueue:
			tg_id, err := h.service.GetTgIdByConvId(ctx, int(message.ConversationID.Int64))
			if err != nil {
				h.logger.Error("Error getting tg id: %s", err)
			}

			if err := sendMessage(tg_id, message, telegramBot); err != nil {
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

			if !update.Message.IsCommand() {
				h.messageToDbQueue <- *update.Message
			} else if err := handleCommands(update, telegramBot); err != nil {
				h.logger.Error("Error sending message to user: %s", err)
			}
		case message := <-h.actionInChat:
			h.sendActionToConvID(ctx, telegramBot, int(message.ConversationID.Int64), "typing")
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

func (h *TelegramHandler) getUpdatesChannel(telegramBot *telegram.Bot) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := telegramBot.Bot.GetUpdatesChan(u)
	return updates
}

func (h *TelegramHandler) sendActionToConvID(ctx context.Context, telegramBot *telegram.Bot, convId int, action_type string) {
	tg_id, err := h.service.GetTgIdByConvId(ctx, convId)
	if err != nil {
		h.logger.Error("Error getting tg id: %s", err)
	}
	action := tgbotapi.NewChatAction(int64(tg_id), "typing")
	telegramBot.Bot.Send(action)
}

func sendMessage(tg_id int, message models.Message, telegramBot *telegram.Bot) error {
	var msg tgbotapi.MessageConfig
	for i := 0; i <= len(message.Answer); i += 4096 {
		if i+4096 < len(message.Answer) {
			msg = tgbotapi.NewMessage(int64(tg_id), message.Answer[i:i+4096])
		} else {
			msg = tgbotapi.NewMessage(int64(tg_id), message.Answer[i:len(message.Answer)])
		}
		msg.ReplyToMessageID = int(message.TgQuestionId.Int64)
		_, err := telegramBot.Bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleCommands(update tgbotapi.Update, telegramBot *telegram.Bot) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		msg.Text = "Hi!"
	default:
		msg.Text = "I don't know that command"
	}
	if _, err := telegramBot.Bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (h *TelegramHandler) startNewConversation(ctx context.Context, id int) error {
	exist, err := h.service.ConvExists(ctx, id)
	if err != nil {
		h.logger.Error("Error checking is user exist: %s", err)
		return err
	}

	if !exist {
		err := h.service.StartNewConversation(ctx, id)
		if err != nil {
			h.logger.Error("Error adding user to db %s", err)
			return err
		}
	}
	return nil
}
