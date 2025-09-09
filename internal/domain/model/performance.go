package model

import (
	"educational-reinforcement-platform/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Period representa o período de tempo para o desempenho
type Period string

const (
	PeriodDaily   Period = "daily"
	PeriodWeekly  Period = "weekly"
	PeriodMonthly Period = "monthly"
	PeriodYearly  Period = "yearly"
)

var (
	ErrInvalidPerformanceData = errors.New("invalid performance data")
	ErrInvalidPeriod          = errors.New("period must be one of: daily, weekly, monthly, yearly")
	ErrInvalidCounter         = errors.New("the counter must be zero or positive")
	ErrSubjectIDEmpty         = errors.New("subject ID cannot be empty")
)

// Performance representa o desempenho do usuário em um determinado período
type Performance struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	SubjectID    string    `json:"subjectId"`
	Period       Period    `json:"period"`
	Correct      int       `json:"correct"`
	Incorrect    int       `json:"incorrect"`
	CalculatedAt time.Time `json:"calculatedAt"`
}

// NewPerformance cria uma nova instância de Performance,
// em caso de erro retorna o erro correspondente
func NewPerformance(userID, subjectID string, period Period, correct, incorrect int) (*Performance, error) {
	const NewPerformanceErrorFmt = "[NewPerformance] ERROR: %w"

	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf(NewPerformanceErrorFmt, err)
	}

	if err := ValidatePerformanceIds(userID, subjectID); err != nil {
		return nil, fmt.Errorf(NewPerformanceErrorFmt, err)
	}

	validPeriod, err := ValidatePeriod(period)
	if err != nil {
		return nil, fmt.Errorf(NewPerformanceErrorFmt, err)
	}

	validCorrect, err := (&Performance{}).ValidatePositiveCounts(correct)
	if err != nil {
		return nil, fmt.Errorf(NewPerformanceErrorFmt, err)
	}

	validIncorrect, err := (&Performance{}).ValidatePositiveCounts(incorrect)
	if err != nil {
		return nil, fmt.Errorf(NewPerformanceErrorFmt, err)
	}

	return &Performance{
		ID:           validId,
		UserID:       userID,
		SubjectID:    subjectID,
		Period:       validPeriod,
		Correct:      validCorrect,
		Incorrect:    validIncorrect,
		CalculatedAt: time.Now(),
	}, nil
}

// ValidatePerformanceIds verifica se os IDs são válidos,
// em caso de erro retorna ErrUserIDEmpty ou ErrSubjectIDEmpty
func ValidatePerformanceIds(userID, subjectID string) error {
	var errs []error
	if len(strings.TrimSpace(userID)) == 0 {
		errs = append(errs, ErrUserIDEmpty)
	}

	if len(strings.TrimSpace(subjectID)) == 0 {
		errs = append(errs, ErrSubjectIDEmpty)
	}

	if len(errs) > 0 {
		return fmt.Errorf(
			"[ValidatePerformanceIds] ERROR: %w, %v",
			ErrInvalidPerformanceData,
			errors.Join(errs...),
		)
	}
	return nil
}

// ValidatePeriod verifica se o período é válido,
// em caso de erro retorna ErrInvalidPeriod
func ValidatePeriod(period Period) (Period, error) {
	switch period {
	case PeriodDaily, PeriodWeekly, PeriodMonthly, PeriodYearly:
		return period, nil
	default:
		return Period(""), fmt.Errorf(
			"[ValidatePeriod] ERROR invalid period(%s): %w",
			period,
			ErrInvalidPeriod,
		)
	}
}

// ValidatePositiveCounts verifica se os contadores de acertos e erros são não negativos,
// em caso de erro retorna ErrInvalidCounter
func (p *Performance) ValidatePositiveCounts(count int) (int, error) {
	if count < 0 {
		return 0, fmt.Errorf(
			"[ValidatePositiveCounts] ERROR invalid count(%d): %w",
			count,
			ErrInvalidCounter,
		)
	}
	return count, nil
}

// UpdateCounts atualiza os contadores de acertos e erros,
// em caso de erro retorna ErrInvalidCounter
func (p *Performance) UpdateCounts(correct, incorrect *int) error {
	if correct != nil {
		validCorrect, err := p.ValidatePositiveCounts(*correct)
		if err != nil {
			return fmt.Errorf("[UpdateCounts] ERROR: %w", err)
		}
		p.Correct = validCorrect
	}

	if incorrect != nil {
		validIncorrect, err := p.ValidatePositiveCounts(*incorrect)
		if err != nil {
			return fmt.Errorf("[UpdateCounts] ERROR: %w", err)
		}
		p.Incorrect = validIncorrect
	}

	p.CalculatedAt = time.Now()
	return nil
}

// GetAccuracy calcula a precisão do desempenho
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
		return fmt.Sprintf("[Performance.String] ERROR: %v", err)
	}
	return string(data)
}
