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
	ErrEmptyQuestionContent   = errors.New("question content cannot be empty")
	ErrQuantityOptions        = errors.New("number of options incompatible with the difficulty")
	ErrInvalidCorrectOptions  = errors.New("there must be exactly one correct option")
	ErrEmptySubjectID         = errors.New("subject ID cannot be empty")
	ErrAddOptionExceedsLimit  = errors.New("cannot add more options than the difficulty allows")
	ErrRemoveOptionBelowLimit = errors.New("cannot have fewer options than the difficulty requires")
	ErrOptionNotFound         = errors.New("option not found")
	ErrChangeDifficulty       = errors.New("current options exceed new difficulty")
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

// NewQuestion cria uma nova instância de Question,
// em caso de erro retorna o erro correspondente
func NewQuestion(subjectID, content string, difficulty Difficulty, options []Option) (*Question, error) {
	const newQuestionErrorFmt = "[NewQuestion] ERROR: %w"

	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf(newQuestionErrorFmt, err)
	}

	validSubjectID, err := ValidateSubjectID(subjectID)
	if err != nil {
		return nil, fmt.Errorf(newQuestionErrorFmt, err)
	}

	validContent, err := ValidateQuestionContent(content)
	if err != nil {
		return nil, fmt.Errorf(newQuestionErrorFmt, err)
	}

	validDifficulty, err := ValidateDifficulty(difficulty)
	if err != nil {
		return nil, fmt.Errorf(newQuestionErrorFmt, err)
	}

	validOptions, err := ValidateOptions(options, validDifficulty)
	if err != nil {
		return nil, fmt.Errorf(newQuestionErrorFmt, err)
	}

	now := time.Now()

	return &Question{
		ID:         validId,
		SubjectID:  validSubjectID,
		Content:    validContent,
		Options:    validOptions,
		Difficulty: validDifficulty,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ValidateSubjectID verifica se o ID do assunto é válido
// em caso de erro retorna ErrEmptySubjectID
func ValidateSubjectID(subjectID string) (string, error) {
	if len(strings.TrimSpace(subjectID)) == 0 {
		return "", fmt.Errorf("[ValidateSubjectID] ERROR: %w", ErrEmptySubjectID)
	}
	return subjectID, nil
}

// ValidateQuestionContent verifica se o conteúdo da pergunta é válido
// em caso de erro retorna ErrEmptyQuestionContent
func ValidateQuestionContent(content string) (string, error) {
	if len(strings.TrimSpace(content)) == 0 {
		return "", fmt.Errorf("[ValidateQuestionContent] ERROR: %w", ErrEmptyQuestionContent)
	}
	return content, nil
}

// ValidateOptions verifica se a lista de opções é válida
// em caso de erro retorna ErrQuantityOptions ou ErrInvalidCorrectOptions
func ValidateOptions(options []Option, difficulty Difficulty) ([]Option, error) {
	if len(options) < int(VeryEasy) || len(options) > int(difficulty) {
		return nil, fmt.Errorf(
			"[ValidateOptions] ERROR: %w (current=%d, max=%d)",
			ErrQuantityOptions,
			len(options),
			difficulty,
		)
	}

	correctCount := 0
	for _, opt := range options {
		if opt.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return nil, fmt.Errorf("[ValidateOptions] ERROR: %w", ErrInvalidCorrectOptions)
	}

	return options, nil
}

// UpdateContent altera o conteúdo da pergunta,
// em caso de erro retorna ErrEmptyQuestionContent
func (q *Question) UpdateContent(newContent string) error {
	validContent, err := ValidateQuestionContent(newContent)
	if err != nil {
		return fmt.Errorf("[UpdateContent] ERROR: %w", err)
	}
	q.Content = validContent
	q.UpdatedAt = time.Now()
	return nil
}

// UpdateDifficulty altera a dificuldade da pergunta,
// em caso de erro retorna ErrChangeDifficulty
func (q *Question) UpdateDifficulty(newDifficulty Difficulty) error {
	validDifficulty, err := ValidateDifficulty(newDifficulty)
	if err != nil {
		return fmt.Errorf("[UpdateDifficulty] ERROR: %w", err)
	}

	if len(q.Options) > int(validDifficulty) {
		return fmt.Errorf(
			"[UpdateDifficulty] ERROR: %w (current options=%d, new difficulty=%d)",
			ErrChangeDifficulty,
			len(q.Options),
			validDifficulty,
		)
	}

	q.Difficulty = validDifficulty
	q.UpdatedAt = time.Now()
	return nil
}

// UpdateOptions altera as opções da pergunta,
// em caso de erro retorna ErrQuantityOptions ou ErrInvalidCorrectOptions
func (q *Question) UpdateOptions(newOptions []Option) error {
	validOptions, err := ValidateOptions(newOptions, q.Difficulty)
	if err != nil {
		return fmt.Errorf("[UpdateOptions] ERROR: %w", err)
	}
	q.Options = validOptions
	q.UpdatedAt = time.Now()
	return nil
}

// AddOption adiciona uma nova opção à pergunta,
// em caso de erro retorna ErrAddOptionExceedsLimit,
// ErrQuantityOptions ou ErrInvalidCorrectOptions
func (q *Question) AddOption(option Option) error {
	if len(q.Options) >= int(q.Difficulty) {
		return fmt.Errorf(
			"[AddOption] ERROR: %w (current=%d, max=%d)",
			ErrAddOptionExceedsLimit,
			len(q.Options), q.Difficulty,
		)
	}

	newOptions := append(q.Options, option)
	_, err := ValidateOptions(newOptions, q.Difficulty)
	if err != nil {
		return fmt.Errorf("[AddOption] ERROR: %w", err)
	}

	q.Options = newOptions
	q.UpdatedAt = time.Now()
	return nil
}

// RemoveOption remove uma opção da pergunta pelo ID,
// em caso de erro retorna ErrRemoveOptionBelowLimit ou ErrOptionNotFound
func (q *Question) RemoveOption(optionID string) error {
	if len(q.Options) <= int(VeryEasy) {
		return fmt.Errorf(
			"[RemoveOption] ERROR: %w (current=%d, min=%d)",
			ErrRemoveOptionBelowLimit,
			len(q.Options),
			VeryEasy,
		)
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
		return fmt.Errorf(
			"[RemoveOption] ERROR Invalid option ID(%s): %w",
			optionID,
			ErrOptionNotFound,
		)
	}

	_, err := ValidateOptions(newOptions, q.Difficulty)
	if err != nil {
		return fmt.Errorf("[RemoveOption] ERROR: %w", err)
	}

	q.Options = newOptions
	q.UpdatedAt = time.Now()
	return nil
}

// SetCorrectOption define qual opção da pergunta deve ser marcada como correta,
// em caso de erro retorna ErrOptionNotFound, ErrQuantityOptions ou ErrInvalidCorrectOptions
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
		return fmt.Errorf(
			"[SetCorrectOption] ERROR Invalid option ID(%s): %w",
			optionID,
			ErrOptionNotFound,
		)
	}

	_, err := ValidateOptions(q.Options, q.Difficulty)
	if err != nil {
		return fmt.Errorf("[SetCorrectOption] ERROR: %w", err)
	}

	q.UpdatedAt = now
	return nil
}

// String retorna a representação em JSON da pergunta
func (q *Question) String() string {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return fmt.Sprintf("[Question.String] ERROR: %v", err)
	}
	return string(data)
}
