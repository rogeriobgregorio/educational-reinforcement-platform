package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Erros específicos do modelo Question
var (
	ErrQuestionIDEmpty      = errors.New("question ID cannot be empty")
	ErrEmptyQuestionContent = errors.New("question content cannot be empty")
)

// Question representa uma pergunta
type Question struct {
	ID         string     `json:"id"`
	SubjectID  string     `json:"subjectId"`
	Content    string     `json:"content"`
	Difficulty Difficulty `json:"difficulty"`
	Options    []Option   `json:"options"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// NewQuestion cria uma nova instância de Question.

// Em caso de erro retorna ValidationError.
func NewQuestion(id, subjectID, content string, difficulty Difficulty, options []Option) (*Question, error) {
	now := time.Now()
	question := &Question{
		ID:         id,
		SubjectID:  subjectID,
		Content:    content,
		Options:    options,
		Difficulty: difficulty,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := question.Validate(); err != nil {
		return nil, err
	}
	return question, nil
}

// Validate verifica se os dados da pergunta são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados
func (q *Question) Validate() error {
	ve := &ValidationError{}

	if strings.TrimSpace(q.ID) == "" {
		ve.Add(ErrQuestionIDEmpty)
	}

	if strings.TrimSpace(q.SubjectID) == "" {
		ve.Add(ErrSubjectIDEmpty)
	}

	if err := validateQuestionContent(q.Content); err != nil {
		ve.Add(err)
	}

	if err := validateDifficulty(q.Difficulty); err != nil {
		ve.Add(err)
	}

	if err := validateOptions(q.Options, q.Difficulty); err != nil {
		ve.Add(err)
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}

// validateQuestionContent verifica se o conteúdo da pergunta é válido.
//
// Em caso de erro retorna ErrEmptyQuestionContent.
func validateQuestionContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return ErrEmptyQuestionContent
	}
	return nil
}

// validateOptions verifica se a lista de opções é válida.
//
// Em caso de erro retorna: ErrQuantityOptions ou ErrInvalidCorrectOptions
func validateOptions(options []Option, difficulty Difficulty) error {
	if len(options) < int(VeryEasy) || len(options) > int(difficulty) {
		return ErrQuantityOptions
	}

	correctCount := 0
	for _, opt := range options {
		if opt.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return ErrInvalidCorrectOptions
	}

	return nil
}

// UpdateContent altera o conteúdo da pergunta.

// Em caso de erro retorna ErrEmptyQuestionContent.
func (q *Question) UpdateContent(newContent string) error {
	if err := validateQuestionContent(newContent); err != nil {
		return err
	}
	q.Content = newContent
	q.UpdatedAt = time.Now()
	return nil
}

// UpdateDifficulty altera a dificuldade da pergunta.
//
// Em caso de erro retorna ErrChangeDifficulty.
func (q *Question) UpdateDifficulty(newDifficulty Difficulty) error {
	if err := validateDifficulty(newDifficulty); err != nil {
		return err
	}

	if len(q.Options) > int(newDifficulty) {
		return ErrChangeDifficulty
	}

	q.Difficulty = newDifficulty
	q.UpdatedAt = time.Now()
	return nil
}

// UpdateOptions altera as opções da pergunta.
//
// Em caso de erro retorna ErrQuantityOptions ou ErrInvalidCorrectOptions.
func (q *Question) UpdateOptions(newOptions []Option) error {
	if err := validateOptions(newOptions, q.Difficulty); err != nil {
		return err
	}
	q.Options = newOptions
	q.UpdatedAt = time.Now()
	return nil
}

// AddOption adiciona uma nova opção à pergunta.
//
// Em caso de erro retorna: ErrAddOptionExceedsLimit, ErrQuantityOptions ou ErrInvalidCorrectOptions.
func (q *Question) AddOption(option Option) error {
	if len(q.Options) >= int(q.Difficulty) {
		return ErrAddOptionExceedsLimit
	}

	newOptions := append(q.Options, option)
	if err := validateOptions(newOptions, q.Difficulty); err != nil {
		return err
	}

	q.Options = newOptions
	q.UpdatedAt = time.Now()
	return nil
}

// RemoveOption remove uma opção da pergunta pelo ID.
//
// Em caso de erro retorna ErrRemoveOptionBelowLimit ou ErrOptionNotFound.
func (q *Question) RemoveOption(optionID string) error {
	if len(q.Options) <= int(VeryEasy) {
		return ErrRemoveOptionBelowLimit
	}

	var newOptions []Option
	var found bool

	for _, opt := range q.Options {
		if opt.ID == optionID {
			found = true
			continue
		}
		newOptions = append(newOptions, opt)
	}

	if !found {
		return ErrOptionNotFound
	}

	if err := validateOptions(newOptions, q.Difficulty); err != nil {
		return err
	}

	q.Options = newOptions
	q.UpdatedAt = time.Now()
	return nil
}

// SetCorrectOption define qual opção da pergunta deve ser marcada como correta.
//
// Em caso de erro retorna: ErrOptionNotFound, ErrQuantityOptions ou ErrInvalidCorrectOptions.
func (q *Question) SetCorrectOption(optionID string) error {
	found := false
	now := time.Now()

	for i := range q.Options {
		if q.Options[i].ID == optionID {
			q.Options[i].IsCorrect = true
			q.Options[i].UpdatedAt = now
			found = true
		} else {
			q.Options[i].IsCorrect = false
			q.Options[i].UpdatedAt = now
		}
	}

	if !found {
		return ErrOptionNotFound
	}

	if err := validateOptions(q.Options, q.Difficulty); err != nil {
		return err
	}

	q.UpdatedAt = now
	return nil
}

// String retorna a representação em JSON da pergunta
func (q *Question) String() string {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.Question.String] ERROR: %v", err)
	}
	return string(data)
}
