package delete

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL string `json:"url"`
}

type URLDeleter interface {
	DeleteURL(url string) error
}

func New(log *slog.Logger, storage URLDeleter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req *Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Info("cant decode json")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err = storage.DeleteURL(req.URL)
		if err != nil {
			log.Info("cant delete url")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		render.JSON(w, r, response.OK())
	})
}
