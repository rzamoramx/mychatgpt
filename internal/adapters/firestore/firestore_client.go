package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"

	"my_chat_gpt/internal/domain"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type FirestoreCLient struct {
	Ctx         context.Context
	MsgsColl    string
	HistoryColl string
	UsersColl   string

	// internal use
	fs *firestore.Client
}

func NewFirestoreClient(msgsColl string, historyColl string, usersColl string) (any, error) {
	if msgsColl == "" || historyColl == "" || usersColl == "" {
		return FirestoreCLient{}, errors.New("messages or History collections name are empty")
	}

	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return FirestoreCLient{}, err
	}

	return &FirestoreCLient{
		Ctx:         ctx,
		MsgsColl:    msgsColl,
		HistoryColl: historyColl,
		UsersColl:   usersColl,
		fs:          client,
	}, nil
}

func (class *FirestoreCLient) DeleteHistory(historyId string) error {
	if historyId == "" {
		return errors.New("history id is empty")
	}

	query := class.fs.Collection(class.MsgsColl).Where("history_id", "==", historyId)

	// Ejecuta la consulta y obtén un iterador que contenga todos los documentos que cumplen con los criterios de la consulta
	docs, err := query.Documents(class.Ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error al obtener los documentos: %v", err)
	}

	// Borra todos los documentos devueltos por la consulta
	for _, doc := range docs {
		if _, err := doc.Ref.Delete(class.Ctx); err != nil {
			fmt.Printf("error al borrar el documento %v: %v", doc.Ref.ID, err)
		}
		fmt.Printf("Documento %v borrado\n", doc.Ref.ID)
	}

	_, err = class.fs.Collection(class.HistoryColl).Doc(historyId).Delete(class.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (class *FirestoreCLient) GetAllHistory() ([]domain.History, error) {
	query := class.fs.Collection(class.HistoryColl)
	docs, err := query.Documents(class.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	history := []domain.History{}
	for _, doc := range docs {
		h := domain.History{}
		err = doc.DataTo(&h)
		if err != nil {
			return history, err
		}
		history = append(history, h)
	}

	return history, nil
}

func (class *FirestoreCLient) SaveHistory(history domain.History) (string, error) {
	if history == (domain.History{}) {
		return "", errors.New("history is empty")
	}

	if history.Id == "" {
		history.Id = uuid.New().String()
	}

	_, err := class.fs.Collection(class.HistoryColl).Doc(history.Id).Set(class.Ctx, history)
	if err != nil {
		return "", err
	}

	return history.Id, nil
}

func (class *FirestoreCLient) SaveMessage(message domain.Message) error {
	if message == (domain.Message{}) {
		return errors.New("message is empty")
	}

	if message.Id == "" {
		message.Id = uuid.New().String()
	}

	_, err := class.fs.Collection(class.MsgsColl).Doc(message.Id).Set(class.Ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (class *FirestoreCLient) GetMessagesByHistoryId(historyId string, order string, limit int) ([]domain.Message, error) {
	if historyId == "" {
		return []domain.Message{}, errors.New("historyId is empty")
	}

	// get by history id and order by timestamp
	messages := []domain.Message{}

	query := class.fs.Collection(class.MsgsColl).Where("history_id", "==", historyId)

	if order == "asc" {
		query = query.OrderBy("timestamp", firestore.Asc)
	} else if order == "desc" {
		query = query.OrderBy("timestamp", firestore.Desc)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	docs, err := query.Documents(class.Ctx).GetAll()
	if err != nil {
		return messages, fmt.Errorf("error on query: %v", err)
	}

	for _, doc := range docs {
		message := domain.Message{}
		err = doc.DataTo(&message)
		if err != nil {
			return messages, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (class *FirestoreCLient) GetUser(username string, pwd string) (domain.User, error) {
	if username == "" || pwd == "" {
		return domain.User{}, errors.New("username or pwd is empty")
	}

	query := class.fs.Collection(class.UsersColl).Where("username", "==", username).Where("pwd", "==", pwd)

	docs, err := query.Documents(class.Ctx).GetAll()
	if err != nil {
		return domain.User{}, fmt.Errorf("error on query: %v", err)
	}

	if len(docs) == 0 {
		return domain.User{}, errors.New("user not found")
	}

	user := domain.User{}
	err = docs[0].DataTo(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
