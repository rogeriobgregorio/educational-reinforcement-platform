package model

import (
	"educational-reinforcement-platform/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidSubjectName = errors.New("subject name cannot be less than 3 characters")

// Subject representa uma disciplina ou matéria
type Subject struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewSubject cria uma nova instância de Subject,
// em caso de erro retorna o erro correspondente
func NewSubject(name string) (*Subject, error) {
	const newSubjectErrorFmt = "[NewSubject] ERROR: %w"

	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf(newSubjectErrorFmt, err)
	}

	validName, err := ValidateSubjectName(name)
	if err != nil {
		return nil, fmt.Errorf(newSubjectErrorFmt, err)
	}

	now := time.Now()

	return &Subject{
		ID:        validId,
		Name:      validName,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ValidateSubjectName verifica se o nome da disciplina é válido,
// em caso de erro retorna ErrInvalidSubjectName
func ValidateSubjectName(name string) (string, error) {
	if len(strings.TrimSpace(name)) < 3 {
		return "", fmt.Errorf(
			"[ValidateSubjectName] ERROR: invalid subject name(%s), %w",
			name,
			ErrInvalidSubjectName,
		)
	}
	return name, nil
}

// UpdateName atualiza o nome da disciplina,
// em caso de erro retorna ErrInvalidSubjectName
func (s *Subject) UpdateName(newName string) error {
	validName, err := ValidateSubjectName(newName)
	if err != nil {
		return fmt.Errorf(
			"[UpdateName] ERROR: invalid subject name(%s), %w",
			newName,
			err,
		)
	}

	s.Name = validName
	s.UpdatedAt = time.Now()
	return nil
}

// String retorna uma representação em JSON da disciplina
func (s *Subject) String() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Sprintf("[Subject.String] ERROR: %v", err)
	}
	return string(data)
}
