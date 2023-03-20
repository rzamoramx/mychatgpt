package services

import (
	"my_chat_gpt/internal/domain"
)

type AiProvider interface {
	GetAnswer(params map[string]string, oldMessages []domain.Message) (string, error)
}
