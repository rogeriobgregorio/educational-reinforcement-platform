package model

import (
	"educational-reinforcement-platform/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrEmptyOptionContent = errors.New("option content cannot be empty")
	ErrEmptyQuestionID    = errors.New("question ID cannot be empty")
)

// Option representa uma opção de resposta para uma pergunta
type Option struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"questionId"`
	Content    string    `json:"content"`
	IsCorrect  bool      `json:"isCorrect"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// NewOption cria uma nova instância de Option
func NewOption(questionID, content string, isCorrect bool) (*Option, error) {
	const newOptionErrorFmt = "[NewOption] ERROR: %w"

	validQuestionID, err := ValidateQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf(newOptionErrorFmt, err)
	}

	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf(newOptionErrorFmt, err)
	}

	validContent, err := ValidateOptionContent(content)
	if err != nil {
		return nil, fmt.Errorf(newOptionErrorFmt, err)
	}

	now := time.Now()

	return &Option{
		ID:         validId,
		QuestionID: validQuestionID,
		Content:    validContent,
		IsCorrect:  isCorrect,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ValidateOptionContent verifica se o conteúdo da opção é válido
func ValidateOptionContent(content string) (string, error) {
	if len(strings.TrimSpace(content)) == 0 {
		return "", fmt.Errorf("[ValidateOptionContent] ERROR: %w", ErrEmptyOptionContent)
	}
	return content, nil
}

// ValidateQuestionID verifica se o ID da pergunta é válido
func ValidateQuestionID(questionID string) (string, error) {
	if len(strings.TrimSpace(questionID)) == 0 {
		return "", fmt.Errorf("[ValidateQuestionID] ERROR: %w", ErrEmptyQuestionID)
	}
	return questionID, nil
}

// ChangeContent altera o conteúdo da opção
func (o *Option) ChangeContent(newContent string) error {
	validContent, err := ValidateOptionContent(newContent)
	if err != nil {
		return fmt.Errorf("[ChangeContent] ERROR: %w", err)
	}
	o.Content = validContent
	o.UpdatedAt = time.Now()
	return nil
}

// String retorna a representação em JSON da opção
func (o *Option) String() string {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return fmt.Sprintf("[Option.String] ERROR: %v", err)
	}
	return string(data)
}
