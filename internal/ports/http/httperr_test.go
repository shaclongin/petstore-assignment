package ports_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tanninio/home-assignment/internal/common"
	ports "github.com/tanninio/home-assignment/internal/ports/http"
)

func TestRespondWithHttpError(t *testing.T) {
	var tests = []struct {
		name     string
		input    error
		wantcode int
		wantbody string
	}{
		{"incorrect-input", common.ErrIncorrectInput, 400, "{\"error\": \"Bad Request\"}"},
		{"unimplemented", common.ErrUnimplemented, 501, "{\"error\": \"Not Implemented\"}"},
		{"unknown", common.ErrUnknown, 500, "{\"error\": \"Internal Server Error\"}"},
		{"not-found", common.ErrNotFound, 404, "{\"error\": \"Not Found\"}"},
		{"already-exists", common.ErrAlreadyExists, 409, "{\"error\": \"Conflict\"}"},
		{"wrapped", fmt.Errorf("wrapped: %w", common.ErrNotFound), 404, "{\"error\": \"Not Found\"}"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()
			ports.HttpRespondWithHttpError(w, r, tt.input)
			require.Equal(t, tt.wantcode, w.Code)
			body, _ := io.ReadAll(w.Body)
			require.JSONEq(t, tt.wantbody, string(body))
		})
	}
}
