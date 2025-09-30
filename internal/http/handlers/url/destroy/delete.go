package destroy

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Magic-B/url-shortener/internal/storage"
	resp "github.com/Magic-B/url-shortener/pkg/http/response"
	"github.com/Magic-B/url-shortener/pkg/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))

			return 
		}

		logger.Info("alias readed", slog.String("alias", alias))

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			
			return
		}
		if err != nil {
			logger.Error("failed to delete url", slg.Error(err))
			render.JSON(w, r, "internal error")

			return 
		}

		logger.Info("url deleted", slog.String("alias", alias))
		render.Status(r, http.StatusOK)
	}
}