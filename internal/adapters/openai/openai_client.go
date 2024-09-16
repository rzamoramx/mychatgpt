package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"my_chat_gpt/configs"
	"my_chat_gpt/internal/domain"
	"net/http"
)

type OpenaiClient struct {
	client http.Client
}

func NewOpenaiClient() any {
	return &OpenaiClient{client: http.Client{}}
}

func (c *OpenaiClient) GetAnswer(params map[string]string, oldMessages []domain.Message) (string, error) {
	role := params["role"]
	prompt := params["prompt"]
	msgs := []Message{}

	// inverse order of oldMessages
	for i := len(oldMessages)/2 - 1; i >= 0; i-- {
		opp := len(oldMessages) - 1 - i
		oldMessages[i], oldMessages[opp] = oldMessages[opp], oldMessages[i]
	}

	for _, msg := range oldMessages {
		msgs = append(msgs, Message{
			Role:    "user",
			Content: msg.Prompt,
		})
		msgs = append(msgs, Message{
			Role:    "system",
			Content: msg.Message,
		})
	}

	//if oldMessages, ok := params["oldMessages"]; ok {

	if role == "" && prompt == "" {
		return "", errors.New("role and prompt are empty")
	}

	if role == "" {
		msgs = append(msgs, Message{
			Role:    "user",
			Content: prompt,
		})
	} else {
		msgs = append(msgs, Message{
			Role:    "system",
			Content: role,
		})
		msgs = append(msgs, Message{
			Role:    "user",
			Content: prompt,
		})
	}

	fmt.Printf("What model we are using: %s\n", configs.OPENAI_MODEL)

	// prepare request
	data := RequestToAOpenAi{
		Model:     configs.OPENAI_MODEL,
		MaxTokens: 3000,
		Msgs:      msgs,
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req, err := http.NewRequest("POST", configs.OPENAI_CHAT_COMPLETIONS_URL, bytes.NewBuffer(json_data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+configs.OPENAI_TOKEN)
	req.Header.Add("Content-Type", "application/json")

	// send request
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// read request
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var result ResponseFromAOpenAi
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(body))

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from openai")
	}

	if result.Choices[0].FisnishReason != "stop" {
		return "", fmt.Errorf("openai did not stop")
	}

	return result.Choices[0].Msg.Content, nil
}
