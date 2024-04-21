package http

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eriktate/splitdecision/tpl"
)

func HandleIndex(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, authenticated := GetSessionFromContext(ctx)
		templ.Handler(tpl.Index(staticPath, r.URL.Path, authenticated)).ServeHTTP(w, r)
	}
}
