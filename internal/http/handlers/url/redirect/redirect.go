package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

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

		resUrl, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
		}
		if err != nil {
			logger.Info("failed to get URL", slg.Error(err))
			render.JSON(w, r, resp.Error("internal error"))
		}

		logger.Info("url readed", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}