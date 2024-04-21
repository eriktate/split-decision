package http

import (
	"net/http"
	"time"

	api "github.com/eriktate/splitdecision"
	"github.com/rs/zerolog/log"
)

func makeSessionCookie(id api.ID, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    api.StringFromID(id),
		Expires:  expiresAt,
		HttpOnly: true,
	}
}

func HandleSignup(hasher api.Hasher, userService api.UserService, sessionService api.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			log.Error().Msg("email and password are both required")
			badRequest(w, "email and password must be provided")
			return
		}

		token, err := hasher.HashString(password)
		if err != nil {
			log.Error().Err(err).Msg("failed to hash password")
			serverError(w, "something went wrong")
			return
		}

		user := api.NewUser{
			Email:    email,
			Token:    token,
			AuthType: api.AuthTypeBasic,
		}

		userID, err := userService.CreateUser(ctx, user)
		if err != nil {
			log.Error().Err(err).Msg("failed to create user")
			serverError(w, "could not create user, please try again")
			return
		}

		session := api.NewSession{
			UserID:    userID,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		sessionID, err := sessionService.CreateSession(ctx, session)
		if err != nil {
			log.Error().Err(err).Msg("failed to create session")
			serverError(w, "user created successfully, but a session couldn't be created. Please login")
			return
		}

		http.SetCookie(w, makeSessionCookie(sessionID, session.ExpiresAt))
		noContent(w)
	}
}

func HandleLogin(hasher api.Hasher, userService api.UserService, sessionService api.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			log.Error().Msg("email and password are both required")
			badRequest(w, "email and password must be provided")
			return
		}

		user, err := userService.GetUserByEmail(ctx, email)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch user by email")
			notFound(w, "could not find user with matching email")
			return
		}

		if !hasher.CompareString(user.Token, password) {
			log.Warn().Str("email", email).Msg("failed login attempt")
			unauthorized(w, "invalid credentials provided")
			return
		}

		session := api.NewSession{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		sessionID, err := sessionService.CreateSession(ctx, session)
		if err != nil {
			log.Error().Err(err).Msg("failed to create session")
			serverError(w, "user created successfully, but a session couldn't be created. Please login")
			return
		}

		http.SetCookie(w, makeSessionCookie(sessionID, session.ExpiresAt))
		noContent(w)
	}
}
