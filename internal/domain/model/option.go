package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Erros específicos do modelo Option
var (
	ErrEmptyOptionContent     = errors.New("option content cannot be empty")
	ErrQuantityOptions        = errors.New("number of options incompatible with the difficulty")
	ErrInvalidCorrectOptions  = errors.New("there must be exactly one correct option")
	ErrAddOptionExceedsLimit  = errors.New("cannot add more options than the difficulty allows")
	ErrRemoveOptionBelowLimit = errors.New("cannot have fewer options than the difficulty requires")
	ErrOptionNotFound         = errors.New("option not found")
	ErrOptionIDEmpty          = errors.New("option ID cannot be empty")
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

// NewOption cria uma nova instância de Option.
//
// Em caso de erro retorna ValidationError.
func NewOption(id, questionID, content string, isCorrect bool) (*Option, error) {
	now := time.Now()
	option := &Option{
		ID:         id,
		QuestionID: questionID,
		Content:    content,
		IsCorrect:  isCorrect,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := option.Validate(); err != nil {
		return nil, err
	}

	return option, nil
}

// Validate verifica se os dados da opção são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func (o *Option) Validate() error {
	ve := &ValidationError{}

	if strings.TrimSpace(o.ID) == "" {
		ve.Add(ErrOptionIDEmpty)
	}

	if strings.TrimSpace(o.QuestionID) == "" {
		ve.Add(ErrQuestionIDEmpty)
	}

	if err := validateOptionContent(o.Content); err != nil {
		ve.Add(ErrEmptyOptionContent)
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// validateOptionContent verifica se o conteúdo da opção é válido.
//
// Em caso de erro retorna ErrEmptyOptionContent.
func validateOptionContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return ErrEmptyOptionContent
	}
	return nil
}

// UpdateContent atualiza o conteúdo da opção.
//
// Em caso de erro retorna ErrEmptyOptionContent.
func (o *Option) UpdateContent(newContent string) error {
	if err := validateOptionContent(newContent); err != nil {
		return err
	}
	o.Content = newContent
	o.UpdatedAt = time.Now()
	return nil
}

// String retorna a representação em JSON da opção
func (o *Option) String() string {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.Option.String] ERROR: %v", err)
	}
	return string(data)
}
