package http

import "context"

// TelegramClient отправляет сообщения в Telegram и обслуживает команды бота.
type TelegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
	// Run запускает long polling для обработки команд (например, /start); блокируется до отмены ctx.
	Run(ctx context.Context)
}
