package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Erros específicos do modelo User
var (
	ErrInvalidName   = errors.New("user name cannot be less than 3 characters")
	ErrInvalidRole   = errors.New("invalid role")
	ErrEmptyRole     = errors.New("role cannot be empty")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrEmptyPassword = errors.New("password hash cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrUserIDEmpty   = errors.New("user ID cannot be empty")
)

// Pattern para validação de email
const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

// Regex compilada para validação de email
var emailRegex = regexp.MustCompile(emailRegexPattern)

// Role define os papéis possíveis para um usuário.
type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

// Status define os status possíveis para um usuário.
type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

// User representa um usuário do sistema.
type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role"`
	Difficulty   Difficulty `json:"difficulty"`
	Status       Status     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// NewUser cria uma nova instância de User.
//
// Em caso de erro retorna ValidationError.
func NewUser(id, name, email, passwordHash string, role Role, difficulty Difficulty) (*User, error) {
	now := time.Now()
	user := &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		Difficulty:   difficulty,
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

// Validate verifica se os dados do usuário são válidos.
//
// Em caso de erro retorna ValidationError que contém todos os erros encontrados.
func (u *User) Validate() error {
	ve := &ValidationError{}

	if strings.TrimSpace(u.ID) == "" {
		ve.Add(ErrUserIDEmpty)
	}

	if err := validateUserName(u.Name); err != nil {
		ve.Add(err)
	}

	if err := validateEmail(u.Email); err != nil {
		ve.Add(err)
	}

	if err := validateRole(u.Role); err != nil {
		ve.Add(err)
	}

	if err := validateDifficulty(u.Difficulty); err != nil {
		ve.Add(err)
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		ve.Add(ErrEmptyPassword)
	}
	if ve.HasErrors() {
		return ve
	}

	return nil
}

// validateUserName verifica se o nome é válido.
//
// Em caso de erro retorna ErrInvalidName.
func validateUserName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return ErrInvalidName
	}
	return nil
}

// validateEmail verifica se o email está em um formato válido.
//
// Em caso de erro retorna: ErrEmptyEmail ou ErrInvalidEmail.
func validateEmail(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return ErrEmptyEmail
	}
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateRole verifica se o papel é válido.
//
// Em caso de erro retorna: ErrInvalidRole ou ErrEmptyRole.
func validateRole(role Role) error {
	switch role {
	case RoleAdmin, RoleUser:
		return nil
	default:
		if role == "" {
			return ErrEmptyRole
		}
		return ErrInvalidRole
	}
}

// UpdateName atualiza o nome do usuário.
//
// Em caso de erro retorna ErrInvalidName.
func (u *User) UpdateName(newName string) error {
	if err := validateUserName(newName); err != nil {
		return err
	}
	u.Name = strings.TrimSpace(newName)
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail atualiza o email do usuário.
//
// Em caso de erro retorna: ErrEmptyEmail ou ErrInvalidEmail.
func (u *User) UpdateEmail(newEmail string) error {
	if err := validateEmail(newEmail); err != nil {
		return err
	}
	u.Email = strings.ToLower(strings.TrimSpace(newEmail))
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateRole atualiza o papel do usuário.
//
// Em caso de erro retorna: ErrInvalidRole ou ErrEmptyRole.
func (u *User) UpdateRole(role Role) error {
	if err := validateRole(role); err != nil {
		return err
	}
	u.Role = role
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateDifficulty atualiza a dificuldade do usuário.
//
// Em caso de erro retorna: ErrInvalidDifficulty ou ErrEmptyDifficulty.
func (u *User) UpdateDifficulty(difficulty Difficulty) error {
	if err := validateDifficulty(difficulty); err != nil {
		return err
	}
	u.Difficulty = difficulty
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
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.Status = StatusInactive
	u.UpdatedAt = time.Now()
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsInactive verifica se o usuário está inativo
func (u *User) IsInactive() bool {
	return u.Status == StatusInactive
}

// String retorna uma representação em JSON do usuário.
//
// Em caso de erro, retorna uma string de erro.
func (u *User) String() string {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Sprintf("[model.User.String] ERROR: %v", err)
	}
	return string(data)
}
