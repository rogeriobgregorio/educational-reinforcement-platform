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
	ErrInvalidAnswerData = errors.New("invalid answer data")
	ErrQuestionIDEmpty   = errors.New("question ID cannot be empty")
	ErrOptionIDEmpty     = errors.New("option ID cannot be empty")
)

// Answer representa uma resposta a uma pergunta
type Answer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	QuestionID string    `json:"questionId"`
	OptionID   string    `json:"optionId"`
	IsCorrect  bool      `json:"isCorrect"`
	CreatedAt  time.Time `json:"createdAt"`
}

// NewAnswer cria uma nova instância de Answer
func NewAnswer(userID, questionID, optionID string, isCorrect bool) (*Answer, error) {
	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("[NewAnswer] ERROR: %w", err)
	}

	if err := ValidateAnswerIDs(userID, questionID, optionID); err != nil {
		return nil, fmt.Errorf("[NewAnswer] ERROR: %w", err)
	}

	return &Answer{
		ID:         validId,
		UserID:     userID,
		QuestionID: questionID,
		OptionID:   optionID,
		IsCorrect:  isCorrect,
		CreatedAt:  time.Now(),
	}, nil
}

// ValidateAnswerIDs verifica se os IDs são válidos
func ValidateAnswerIDs(userID, questionID, optionID string) error {
	var errs []error

	if len(strings.TrimSpace(userID)) == 0 {
		errs = append(errs, ErrUserIDEmpty)
	}

	if len(strings.TrimSpace(questionID)) == 0 {
		errs = append(errs, ErrQuestionIDEmpty)
	}

	if len(strings.TrimSpace(optionID)) == 0 {
		errs = append(errs, ErrOptionIDEmpty)
	}

	if len(errs) > 0 {
		return fmt.Errorf(
			"[ValidateAnswerIds] ERROR: %w, %v",
			ErrInvalidAnswerData,
			errors.Join(errs...),
		)
	}
	return nil
}

// String retorna uma representação em JSON da resposta
func (a *Answer) String() string {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Sprintf("[Answer.String] ERROR: %v", err)
	}
	return string(data)
}
