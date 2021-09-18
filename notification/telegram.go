package notification

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alesanfra/ground-control/conf"
	"github.com/alesanfra/ground-control/scanner"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Server rest api for the agent
type Server struct {
	devices scanner.DeviceMap
	config  conf.TelegramConfig
	debug   bool
}

func NewNotificationService(devices scanner.DeviceMap, config conf.TelegramConfig) *Server {
	return &Server{devices: devices, config: config, debug: false}
}

func (s *Server) Name() string {
	return "Telegram"
}

func (s *Server) Run(ctx context.Context) error {

	bot, err := telegram.NewBotAPI(s.config.Token)
	if err != nil {
		log.Printf("Error while initializing telegram api: %v", err)
	}

	bot.Debug = s.debug
	updates, err := bot.GetUpdatesChan(telegram.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: int(s.config.Timeout.Duration),
	})

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Chat.ID != int64(s.config.ChatId) { // ignore any unauthorized users
			log.Printf(
				"Message from unauthorized user %s (chat %d)",
				update.Message.From.UserName,
				update.Message.MessageID,
			)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		command := update.Message.Text
		var reply string

		switch command {
		case "devices":
			devices, _ := json.Marshal(s.devices.AsList())
			reply = string(devices)
		case "speed":
			reply = "bb"
		default:
			reply = "Command not recognized, available commands: devices, speed"
		}
		msg := telegram.NewMessage(update.Message.Chat.ID, reply)
		//msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	}

	return nil
}
