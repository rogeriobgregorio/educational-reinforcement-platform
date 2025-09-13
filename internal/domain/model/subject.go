package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Erros específicos do modelo Subject
var (
	ErrInvalidSubjectName = errors.New("subject name cannot be less than 3 characters")
	ErrSubjectIDEmpty     = errors.New("subject ID cannot be empty")
)

// Subject representa uma disciplina ou matéria
type Subject struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewSubject cria uma nova instância de Subject.
//
// Em caso de erro retorna ValidationError.
func NewSubject(id, name string) (*Subject, error) {
	now := time.Now()
	subject := &Subject{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := subject.Validate(); err != nil {
		return nil, err
	}
	return subject, nil
}

// Validate verifica se os dados da disciplina são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func (s *Subject) Validate() error {
	ve := &ValidationError{}
	
	if strings.TrimSpace(s.ID) == "" {
		ve.Add(ErrSubjectIDEmpty)
	}

	if err := validateSubjectName(s.Name); err != nil {
		ve.Add(err)
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}

// validateSubjectName verifica se o nome da disciplina é válido.
//
// Em caso de erro retorna ErrInvalidSubjectName.
func validateSubjectName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return ErrInvalidSubjectName
	}
	return nil
}

// UpdateName atualiza o nome da disciplina.
//
// Em caso de erro retorna ErrInvalidSubjectName.
func (s *Subject) UpdateName(newName string) error {
	if err := validateSubjectName(newName); err != nil {
		return err
	}
	s.Name = newName
	s.UpdatedAt = time.Now()
	return nil
}

// String retorna uma representação em JSON da disciplina
func (s *Subject) String() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.Subject.String] ERROR: %v", err)
	}
	return string(data)
}
