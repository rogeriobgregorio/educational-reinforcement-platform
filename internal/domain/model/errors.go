package model

import (
	"strings"
)

// ValidationError contém múltiplos erros de validação
type ValidationError struct {
	Errors []error
}

// Error implementa a interface error para ValidationError
func (e *ValidationError) Error() string {
	if !e.HasErrors() {
		return "validation failed"
	}

	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

// Unwrap permite desembrulhar os erros contidos em ValidationError
func (e *ValidationError) Unwrap() []error {
	return e.Errors
}

// Add adiciona um erro à lista de erros de validação
func (e *ValidationError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors verifica se há erros de validação
func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}
