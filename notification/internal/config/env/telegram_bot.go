package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	telegramBotTokenEnvName  = "TELEGRAM_BOT_TOKEN" //nolint:gosec // G101: это имя переменной окружения, а не сам токен
	telegramBotChatIDEnvName = "TELEGRAM_BOT_CHAT_ID"
)

type telegramBotConfig struct {
	token  string
	chatID string
}

func NewTelegramBotConfig() (*telegramBotConfig, error) {
	token := os.Getenv(telegramBotTokenEnvName)
	if len(token) == 0 {
		return nil, errors.New("telegram bot token is required")
	}

	chatID := os.Getenv(telegramBotChatIDEnvName)
	if len(chatID) == 0 {
		return nil, errors.New("telegram bot chat id is required")
	}

	return &telegramBotConfig{
		token:  token,
		chatID: chatID,
	}, nil
}

func (cfg *telegramBotConfig) Token() string {
	return cfg.token
}

func (cfg *telegramBotConfig) ChatID() string {
	return cfg.chatID
}
