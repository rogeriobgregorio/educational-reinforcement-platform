package model

import (
	"educational-reinforcement-platform/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidName  = errors.New("name cannot be less than 3 characters")
	ErrInvalidRole  = errors.New("invalid role")
	ErrEmptyRole    = errors.New("role cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")
	ErrEmptyEmail   = errors.New("email cannot be empty")
	ErrUserIDEmpty  = errors.New("user ID cannot be empty")
)

// Role define os papéis possíveis para um usuário
type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

// User representa um usuário do sistema
type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role"`
	Difficulty   Difficulty `json:"difficulty"`
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// NewUser cria um novo usuário com os dados fornecidos,
// em caso de erro retorna o erro correspondente
func NewUser(name, email, passwordHash string, role Role, difficulty Difficulty) (*User, error) {
	const newUserErrorFmt = "[NewUser] ERROR: %w"

	validId, err := pkg.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf(newUserErrorFmt, err)
	}

	validName, err := ValidateUserName(name)
	if err != nil {
		return nil, fmt.Errorf(newUserErrorFmt, err)
	}

	validEmail, err := ValidateEmail(email)
	if err != nil {
		return nil, fmt.Errorf(newUserErrorFmt, err)
	}

	validRole, err := ValidateRole(string(role))
	if err != nil {
		return nil, fmt.Errorf(newUserErrorFmt, err)
	}

	validDifficulty, err := ValidateDifficulty(difficulty)
	if err != nil {
		return nil, fmt.Errorf(newUserErrorFmt, err)
	}

	now := time.Now()

	return &User{
		ID:           validId,
		Name:         validName,
		Email:        validEmail,
		PasswordHash: passwordHash,
		Role:         validRole,
		Difficulty:   validDifficulty,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ValidateUserName verifica se o nome é válido,
// em caso de erro retorna ErrInvalidName
func ValidateUserName(name string) (string, error) {
	if len(strings.TrimSpace(name)) < 3 {
		return "", fmt.Errorf(
			"[ValidateUserName] ERROR: invalid name(%s), %w",
			name,
			ErrInvalidName,
		)
	}
	return name, nil
}

// ValidateEmail verifica se o email está em um formato válido,
// em caso de erro retorna ErrEmptyEmail ou ErrInvalidEmail
func ValidateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if len(email) == 0 {
		return "", fmt.Errorf(
			"[ValidateEmail] ERROR: invalid email(%s), %w",
			email,
			ErrEmptyEmail,
		)
	}

	regexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !regexp.MatchString(email) {
		return "", fmt.Errorf(
			"[ValidateEmail] ERROR: invalid email(%s), %w",
			email,
			ErrInvalidEmail,
		)
	}
	return email, nil
}

// ValidateRole verifica se o papel é válido,
// em caso de erro retorna ErrInvalidRole ou ErrEmptyRole
func ValidateRole(userRole string) (Role, error) {
	userRole = strings.TrimSpace(strings.ToUpper(userRole))
	if len(strings.TrimSpace(userRole)) == 0 {
		return Role(""), fmt.Errorf(
			"[ValidateRole] ERROR: invalid role(%s), %w",
			userRole,
			ErrEmptyRole,
		)
	}

	role := Role(userRole)

	switch role {
	case RoleAdmin, RoleUser:
		return role, nil
	default:
		return Role(""), fmt.Errorf(
			"[ValidateRole] ERROR: invalid role(%s), %w",
			userRole,
			ErrInvalidRole,
		)
	}
}

// UpdateName atualiza o nome do usuário - em caso de erro retorna ErrInvalidName
func (u *User) UpdateName(name string) error {
	validName, err := ValidateUserName(name)
	if err != nil {
		return fmt.Errorf("[UpdateName] ERROR: %w", err)
	}
	u.Name = validName
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail atualiza o email do usuário - em caso de erro retorna ErrEmptyEmail ou ErrInvalidEmail
func (u *User) UpdateEmail(email string) error {
	validEmail, err := ValidateEmail(email)
	if err != nil {
		return fmt.Errorf("[UpdateEmail] ERROR: %w", err)
	}
	u.Email = validEmail
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateRole atualiza o papel do usuário - em caso de erro retorna ErrInvalidRole ou ErrEmptyRole
func (u *User) UpdateRole(role Role) error {
	validRole, err := ValidateRole(string(role))
	if err != nil {
		return fmt.Errorf("[UpdateRole] ERROR: %w", err)
	}
	u.Role = validRole
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateDifficulty atualiza a dificuldade do usuário - em caso de erro retorna ErrInvalidDifficulty
func (u *User) UpdateDifficulty(difficulty Difficulty) error {
	validDifficulty, err := ValidateDifficulty(difficulty)
	if err != nil {
		return fmt.Errorf("[UpdateDifficulty] ERROR: %w", err)
	}
	u.Difficulty = validDifficulty
	u.UpdatedAt = time.Now()
	return nil
}

// IsAdmin verifica se o usuário é um administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser verifica se o usuário é um usuário comum
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

// Activate ativa o usuário
func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.Active
}

// IsInactive verifica se o usuário está inativo
func (u *User) IsInactive() bool {
	return !u.Active
}

// String retorna uma representação em JSON do usuário
func (u *User) String() string {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Sprintf("[User.String] ERROR: %v", err)
	}
	return string(data)
}
