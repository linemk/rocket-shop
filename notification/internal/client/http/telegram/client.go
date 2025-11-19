package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const telegramAPIURL = "https://api.telegram.org"

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Client struct {
	token  string
	chatID string
	client *http.Client
	logger Logger
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func NewClient(token, chatID string, logger Logger) *Client {
	return &Client{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
		logger: logger,
	}
}

func (c *Client) SendMessage(ctx context.Context, text string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIURL, c.token)

	reqBody := sendMessageRequest{
		ChatID:    c.chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error(ctx, "Failed to marshal telegram request", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		c.logger.Error(ctx, "Failed to create telegram request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(ctx, "Failed to send telegram message", zap.Error(err))
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error(ctx, "Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(ctx, "Telegram API returned non-200 status", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	c.logger.Info(ctx, "Telegram message sent successfully")
	return nil
}
