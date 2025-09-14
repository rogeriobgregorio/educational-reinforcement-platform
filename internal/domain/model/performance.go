package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Erros específicos do modelo Performance
var (
	ErrPerformanceIDEmpty     = errors.New("performance ID cannot be empty")
	ErrInvalidPerformanceData = errors.New("invalid performance data")
	ErrInvalidPeriod          = errors.New("period must be one of: daily, weekly, monthly, yearly")
	ErrInvalidCounter         = errors.New("the counter must be zero or positive")
)

// Period representa o período de tempo para o desempenho.
type Period string

const (
	PeriodDaily   Period = "daily"
	PeriodWeekly  Period = "weekly"
	PeriodMonthly Period = "monthly"
	PeriodYearly  Period = "yearly"
)

// Performance representa o desempenho do usuário em um determinado período.
type Performance struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	SubjectID    string    `json:"subjectId"`
	Period       Period    `json:"period"`
	Correct      int       `json:"correct"`
	Incorrect    int       `json:"incorrect"`
	CalculatedAt time.Time `json:"calculatedAt"`
}

// NewPerformance cria uma nova instância de Performance.
//
// Em caso de erro retorna ValidationError.
func NewPerformance(id, userID, subjectID string, period Period, correct, incorrect int) (*Performance, error) {
	performance := &Performance{
		ID:           id,
		UserID:       userID,
		SubjectID:    subjectID,
		Period:       period,
		Correct:      correct,
		Incorrect:    incorrect,
		CalculatedAt: time.Now(),
	}

	if err := performance.Validate(); err != nil {
		return nil, err
	}
	return performance, nil
}

// Validate verifica se os dados do desempenho são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func (p *Performance) Validate() error {
	ve := &ValidationError{}

	if strings.TrimSpace(p.ID) == "" {
		ve.Add(ErrPerformanceIDEmpty)
	}

	if strings.TrimSpace(p.UserID) == "" {
		ve.Add(ErrUserIDEmpty)
	}

	if strings.TrimSpace(p.SubjectID) == "" {
		ve.Add(ErrSubjectIDEmpty)
	}

	if err := validatePeriod(p.Period); err != nil {
		ve.Add(err)
	}

	if p.Correct < 0 {
		ve.Add(ErrInvalidCounter)
	}

	if p.Incorrect < 0 {
		ve.Add(ErrInvalidCounter)
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// validatePeriod verifica se o período é válido.
//
// Em caso de erro retorna ErrInvalidPeriod
func validatePeriod(period Period) error {
	switch period {
	case PeriodDaily, PeriodWeekly, PeriodMonthly, PeriodYearly:
		return nil
	default:
		return ErrInvalidPeriod
	}
}

// UpdateCorrect incrementa o contador de acertos em 1.
func (p *Performance) UpdateCorrect() error {
	p.Correct++
	p.CalculatedAt = time.Now()
	return nil
}

// UpdateIncorrect incrementa o contador de erros em 1.
func (p *Performance) UpdateIncorrect() error {
	p.Incorrect++
	p.CalculatedAt = time.Now()
	return nil
}

// ResetCounters zera os contadores de acertos e erros.
func (p *Performance) ResetCounters() error {
	p.Correct = 0
	p.Incorrect = 0
	p.CalculatedAt = time.Now()
	return nil
}

// GetAccuracy calcula a precisão do desempenho.
func (p *Performance) GetAccuracy() float64 {
	total := p.Correct + p.Incorrect
	if total == 0 {
		return 0.0
	}
	return (float64(p.Correct) / float64(total)) * 100
}

// GetTotalQuestions retorna o total de perguntas respondidas
func (p *Performance) GetTotalQuestions() int {
	return p.Correct + p.Incorrect
}

// String retorna uma representação em JSON do desempenho
func (p *Performance) String() string {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.Performance.String] ERROR: %v", err)
	}
	return string(data)
}
