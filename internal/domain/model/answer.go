package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Erros específicos do modelo Answer
var (
	ErrAnswerIDEmpty = errors.New("answer ID cannot be empty")
)

// Answer representa uma resposta a uma pergunta
type Answer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	QuestionID string    `json:"questionId"`
	OptionID   string    `json:"optionId"`
	IsCorrect  bool      `json:"isCorrect"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// NewAnswer cria uma nova instância de Answer.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func NewAnswer(id, userID, questionID, optionID string, isCorrect bool) (*Answer, error) {
	now := time.Now()
	answer := &Answer{
		ID:         id,
		UserID:     userID,
		QuestionID: questionID,
		OptionID:   optionID,
		IsCorrect:  isCorrect,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := answer.Validate(); err != nil {
		return nil, err
	}
	return answer, nil
}

// Validate verifica se os dados da resposta são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func (a *Answer) Validate() error {
	ve := &ValidationError{}

	if strings.TrimSpace(a.ID) == "" {
		ve.Add(ErrAnswerIDEmpty)
	}

	if strings.TrimSpace(a.UserID) == "" {
		ve.Add(ErrUserIDEmpty)
	}

	if strings.TrimSpace(a.QuestionID) == "" {
		ve.Add(ErrQuestionIDEmpty)
	}

	if strings.TrimSpace(a.OptionID) == "" {
		ve.Add(ErrOptionIDEmpty)
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// String retorna uma representação em JSON da resposta.
//
// Em caso de erro retorna uma string de erro.
func (a *Answer) String() string {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.Answer.String] ERROR: %v", err)
	}
	return string(data)
}
