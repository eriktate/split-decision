package api

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type AuthType string

const (
	AuthTypeBasic  = AuthType("basic")
	AuthTypeGoogle = AuthType("google")
	AuthTypeNone   = AuthType("none")
)

func ParseAuthType(input string) (AuthType, error) {
	authType := AuthType(strings.ToLower(input))
	switch authType {
	case AuthTypeBasic, AuthTypeGoogle:
		return authType, nil
	default:
		return AuthTypeNone, fmt.Errorf("invalid AuthTupe '%s'", authType)
	}
}

func (at AuthType) Value() (driver.Value, error) {
	return string(at), nil
}

func (at *AuthType) Scan(val any) error {
	if val == nil {
		return nil
	}

	if stringVal, err := driver.String.ConvertValue(val); err == nil {
		if v, ok := stringVal.(string); ok {
			var err error
			*at, err = ParseAuthType(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type User struct {
	ID       ID       `json:"id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	AuthType AuthType `json:"authType"`
	Token    []byte   `json:"token"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewUser struct {
	Email    string   `json:"email"`
	AuthType AuthType `json:"authType"`
	Token    []byte   `json:"token"`
}

type UserService interface {
	CreateUser(ctx context.Context, user NewUser) (ID, error)
	GetUser(ctx context.Context, id ID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	DeleteUser(ctx context.Context, id ID) error
}

type Session struct {
	ID        ID        `json:"id"`
	UserID    ID        `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type NewSession struct {
	UserID    ID        `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SessionService interface {
	CreateSession(ctx context.Context, session NewSession) (ID, error)
	GetSession(ctx context.Context, id ID) (Session, error)
	GetUserSession(ctx context.Context, userID ID) (Session, error)
	ExtendSession(ctx context.Context, id ID, expiresAt time.Time) error
	DeleteSession(ctx context.Context, id ID) error
	DeleteUserSessions(ctx context.Context, userID ID) error
}

func (s Session) Valid() bool {
	return s.ExpiresAt.After(time.Now())
}
