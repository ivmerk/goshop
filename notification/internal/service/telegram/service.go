package telegram

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/notification/internal/client/http"
	"github.com/ivmerk/goshop/notification/internal/model"
	"github.com/ivmerk/goshop/platform/pkg/logger"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type service struct {
	client http.TelegramClient
	chatID int64

	paidTemplate      *template.Template
	assembledTemplate *template.Template
}

func NewService(client http.TelegramClient, chatID int64) (*service, error) {
	paid, err := template.ParseFS(templatesFS, "templates/paid_notification.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse paid template: %w", err)
	}

	assembled, err := template.ParseFS(templatesFS, "templates/assembled_notification.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse assembled template: %w", err)
	}

	return &service{
		client:            client,
		chatID:            chatID,
		paidTemplate:      paid,
		assembledTemplate: assembled,
	}, nil
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	return s.send(ctx, s.paidTemplate, event, event.OrderUUID)
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	return s.send(ctx, s.assembledTemplate, event, event.OrderUUID)
}

func (s *service) send(ctx context.Context, tmpl *template.Template, data any, orderUUID string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if err := s.client.SendMessage(ctx, s.chatID, buf.String()); err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	logger.Info(ctx, "Telegram notification sent", zap.String("order_uuid", orderUUID))

	return nil
}
