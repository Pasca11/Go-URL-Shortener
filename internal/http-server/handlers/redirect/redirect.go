package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	st "url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, storage URLGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("empty alias")
			render.JSON(w, r, response.Error("invalid alias"))

			return
		}

		url, err := storage.GetURL(alias)
		if errors.Is(err, st.ErrUrlNotFound) {
			log.Info("URL not found", "alias", alias)

			w.WriteHeader(404)

			render.JSON(w, r, response.Error("not found"))
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	})
}
