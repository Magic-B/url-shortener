package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Magic-B/url-shortener/internal/storage"
	resp "github.com/Magic-B/url-shortener/pkg/http/response"
	"github.com/Magic-B/url-shortener/pkg/logger/slg"
	"github.com/Magic-B/url-shortener/pkg/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

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

		logger.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)

			logger.Error("invalid request", slg.Error(err))

			render.JSON(w, r, resp.ValidationErrors(validateErrs))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLExist) {
			errStr := "url already exists"

			logger.Info(errStr, slog.String("url", req.URL))
			render.JSON(w, r, resp.Error(errStr))

			return 
		}

		if err != nil {
			errStr := "failed to add url"

			logger.Error(errStr, slg.Error(err))
			render.JSON(w, r, resp.Error(errStr))

			return 
		}

		logger.Info("url added ", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias: alias,
		})
	}
}
