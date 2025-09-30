package save

import (
	"log/slog"
	"net/http"

	resp "github.com/Magic-B/url-shortener/pkg/http/response"
	"github.com/Magic-B/url-shortener/pkg/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL   string `json:"url" validate:"required,url" `
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias  string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			const errStr = "failed to decode request body"
			logger.Error(errStr, slg.Error(err))
			render.JSON(w, r, resp.Error(errStr))

			return 
		}
	}
}
