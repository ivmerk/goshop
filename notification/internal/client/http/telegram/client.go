package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/ivmerk/goshop/notification/internal/client/http"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

const welcomeMessage = "Привет! Я бот GoShop 🛍\nСюда будут приходить уведомления об оплате и сборке ваших заказов."

type client struct {
	bot *bot.Bot
}

var _ http.TelegramClient = (*client)(nil)

func NewClient(token string) (*client, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	return &client{bot: b}, nil
}

func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})

	return err
}

func (c *client) Run(ctx context.Context) {
	c.bot.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeMessage,
	}); err != nil {
		logger.Error(ctx, "failed to reply to /start", zap.Error(err))
	}
}
