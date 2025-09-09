package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidDifficulty = errors.New("difficulty must be between VeryEasy(2) and VeryHard(6)")

// Difficulty representa os níveis de dificuldade disponíveis
type Difficulty int

const (
	VeryEasy Difficulty = iota + 2
	Easy
	Medium
	Hard
	VeryHard
)

// ValidateDifficulty verifica se o nível de dificuldade é válido,
// em caso de erro retorna ErrInvalidDifficulty
func ValidateDifficulty(difficulty Difficulty) (Difficulty, error) {
	if difficulty < VeryEasy || difficulty > VeryHard {
		return Difficulty(0), fmt.Errorf(
			"[ValidateDifficulty] ERROR invalid difficulty(%d): %w",
			difficulty,
			ErrInvalidDifficulty,
		)
	}
	return difficulty, nil
}

// String retorna a representação em string do nível de dificuldade
func (d Difficulty) String() string {
	switch d {
	case VeryEasy:
		return "Very Easy"
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case VeryHard:
		return "Very Hard"
	default:
		return "Unknown"
	}
}

// MarshalJSON implementa a interface json.Marshaler para customizar a
// serialização do nível de dificuldade
func (d *Difficulty) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implementa a interface json.Unmarshaler para customizar a
// desserialização do nível de dificuldade,
// em caso de erro retorna ErrInvalidDifficulty
func (d *Difficulty) UnmarshalJSON(data []byte) error {
	var difficultyStr string
	if err := json.Unmarshal(data, &difficultyStr); err != nil {
		return fmt.Errorf("[UnmarshalJSON] ERROR unmarshaling difficulty: %w", err)
	}

	var difficulty Difficulty
	switch difficultyStr {
	case "Very Easy":
		difficulty = VeryEasy
	case "Easy":
		difficulty = Easy
	case "Medium":
		difficulty = Medium
	case "Hard":
		difficulty = Hard
	case "Very Hard":
		difficulty = VeryHard
	default:
		return fmt.Errorf(
			"[UnmarshalJSON] ERROR invalid difficulty string(%s): %w",
			difficultyStr,
			ErrInvalidDifficulty,
		)
	}

	*d = difficulty
	return nil
}

// ToInt converte o nível de dificuldade para um valor inteiro
func (d Difficulty) ToInt() int {
	return int(d)
}

// FromInt converte um valor inteiro para o nível de dificuldade correspondente,
// em caso de erro retorna ErrInvalidDifficulty
func FromInt(value int) (Difficulty, error) {
	difficulty := Difficulty(value)
	return ValidateDifficulty(difficulty)
}
