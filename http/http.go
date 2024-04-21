package http

import (
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
	api "github.com/eriktate/splitdecision"
	"github.com/eriktate/splitdecision/tpl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Addr           string
	StaticDir      string
	Hasher         api.Hasher
	UserService    api.UserService
	SessionService api.SessionService
}

type Middleware func(http.Handler) http.Handler

func getRoutes(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// handle serving static files
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Get("/", fileServer.ServeHTTP)

	// webhook routes
	r.Post("/signup", HandleSignup(cfg.Hasher, cfg.UserService, cfg.SessionService))
	r.Post("/login", HandleLogin(cfg.Hasher, cfg.UserService, cfg.SessionService))
	r.Get("/signup", templ.Handler(tpl.Signup()).ServeHTTP)
	r.With(IncludeSession(cfg.SessionService)).Get("/home", HandleHomeView())

	return r
}

func Serve(cfg Config) error {
	return http.ListenAndServe(cfg.Addr, getRoutes(cfg))
}

// response writers
func respond(w http.ResponseWriter, payload []byte) {
	if _, err := w.Write(payload); err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func respondWithJSON(w http.ResponseWriter, payload any, status int) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload into response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	respond(w, data)
}

func okJSON(w http.ResponseWriter, payload any) {
	respondWithJSON(w, payload, http.StatusOK)
}

func createdJSON(w http.ResponseWriter, payload any) {
	respondWithJSON(w, payload, http.StatusCreated)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	respond(w, nil)
}

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	respond(w, []byte(msg))
}

func serverError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	respond(w, []byte(msg))
}

func notFound(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusNotFound)
	respond(w, []byte(msg))
}

func forbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	respond(w, []byte(msg))
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	respond(w, []byte(msg))
}
