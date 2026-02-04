package ports

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tanninio/home-assignment/internal/common"
)

type HttpErrorResponse struct {
	Slug       string `json:"error"`
	statusCode int
}

type HttpError struct {
	err        error
	msg        string
	statusCode int
}

func (e HttpError) Error() string {
	return e.msg
}

func (e HttpError) Unwrap() error {
	return e.err
}

func (e HttpError) StatusCode() int {
	return e.statusCode
}

var errToStatusCode = map[error]int{
	common.ErrIncorrectInput: http.StatusBadRequest,
	common.ErrUnimplemented:  http.StatusNotImplemented,
	common.ErrUnknown:        http.StatusInternalServerError,
	common.ErrNotFound:       http.StatusNotFound,
	common.ErrAlreadyExists:  http.StatusConflict,
}

func NewHttpError(err error) HttpError {
	statusCode := http.StatusInternalServerError
	for k, v := range errToStatusCode {
		if errors.Is(err, k) {
			statusCode = v
			break
		}
	}
	return HttpError{err: err, msg: err.Error(), statusCode: statusCode}
}

func HttpRespondWithHttpError(w http.ResponseWriter, r *http.Request, err error) {
	httperr := NewHttpError(err)
	statusCode := httperr.StatusCode()
	slug := http.StatusText(statusCode)

	logrus.WithError(err).Error(slug)
	respondWithErrorResponse(w, r, HttpErrorResponse{slug, statusCode})
}

func respondWithErrorResponse(w http.ResponseWriter, r *http.Request, resp HttpErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
