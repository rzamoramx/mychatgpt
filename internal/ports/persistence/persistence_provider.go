package persistence

import (
	"my_chat_gpt/internal/domain"
)

type PersistenceProvider interface {
	GetAllHistory() ([]domain.History, error)
	SaveHistory(history domain.History) (string, error)
	SaveMessage(message domain.Message) error
	GetMessagesByHistoryId(historyId string, order string, limit int) ([]domain.Message, error)
	DeleteHistory(historyId string) error
	GetUser(username string, pwd string) (domain.User, error)
}
