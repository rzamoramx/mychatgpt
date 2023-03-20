package application

import (
	"errors"
	"my_chat_gpt/internal/domain"
	"my_chat_gpt/internal/ports/persistence"
	"my_chat_gpt/internal/ports/services"
	"time"
)

type MyChatGptApp struct {
	aiProvider services.AiProvider
	persister  persistence.PersistenceProvider
}

func NewMyChatGptApp(provider services.AiProvider, persister1 persistence.PersistenceProvider) (MyChatGptApp, error) {
	if provider == nil {
		return MyChatGptApp{}, errors.New("AiProvider is nil")
	}

	if persister1 == nil {
		return MyChatGptApp{}, errors.New("PersistenceProvider is nil")
	}

	return MyChatGptApp{aiProvider: provider, persister: persister1}, nil
}

func (class *MyChatGptApp) DeleteHistory(historyId string) error {
	return class.persister.DeleteHistory(historyId)
}

func (class *MyChatGptApp) RetrieveAllHistories() ([]domain.History, error) {
	return class.persister.GetAllHistory()
}

func (class *MyChatGptApp) GetHistory(historyId string) ([]domain.Message, error) {
	return class.persister.GetMessagesByHistoryId(historyId, "asc", 0)
}

func (class *MyChatGptApp) NewHistory(name string) (string, error) {
	newHistory := domain.History{
		UserId: "1",
		Name:   name,
	}

	return class.persister.SaveHistory(newHistory)
}

func (class *MyChatGptApp) ProcessMessage(historyId string, message string) (string, error) {
	params := map[string]string{
		"role":       "",
		"prompt":     message,
		"max_tokens": "10",
	}

	oldMessages, err := class.persister.GetMessagesByHistoryId(historyId, "desc", 10)
	if err != nil {
		return "", err
	}

	result, err := class.aiProvider.GetAnswer(params, oldMessages)
	if err != nil {
		return "", err
	}

	// persist message
	msgToPersist := domain.Message{
		HistoryId: historyId,
		Role:      "",
		Prompt:    message,
		Message:   result,
		Timestamp: time.Now().Unix(),
	}
	err = class.persister.SaveMessage(msgToPersist)

	return result, err
}

func (class *MyChatGptApp) GetHistoryMessages(historyId string) ([]map[string]string, error) {
	msgs, err := class.persister.GetMessagesByHistoryId(historyId, "asc", 0)
	if err != nil {
		if err.Error() == "no more items in iterator" {
			return []map[string]string{}, nil
		} else {
			return nil, err
		}
	}

	// for every message get prompt as user and message as system and put in map with the keys: from = user or system and text = prompt or message
	toReturn := []map[string]string{}
	for _, msg := range msgs {
		toReturn = append(toReturn, map[string]string{
			"from": "user",
			"text": msg.Prompt,
		})
		toReturn = append(toReturn, map[string]string{
			"from": "system",
			"text": msg.Message,
		})
	}

	return toReturn, nil
}
