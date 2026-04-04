package helpers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AddInfo(r *chi.Mux) {
	r.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World!"))
	})
}
