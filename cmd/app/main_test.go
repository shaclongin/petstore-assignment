package main_test

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"net/http/httptest"
	"runtime/debug"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	main "github.com/tanninio/home-assignment/cmd/app"
	"github.com/tanninio/home-assignment/internal/app"
)

func TestGoMod(t *testing.T) {
	const wanted = "10af06b9ea8760aae829938b5178e47a"
	bi, ok := debug.ReadBuildInfo()
	require.True(t, ok, "failed to read build info")
	sort.Slice(bi.Deps, func(i, j int) bool {
		return bi.Deps[i].Path < bi.Deps[j].Path
	})
	h := fnv.New128()
	for _, x := range bi.Deps {
		h.Write([]byte(x.Sum))
	}
	got := hex.EncodeToString(h.Sum(nil))
	require.Equal(t, wanted, got, "go.mod must not change")
}

func TestEndToEnd(t *testing.T) {
	existingpet := app.Pet{Id: 1337, Name: "Leet Dog"}
	existingreq := "{\"name\": \"Leet Dog\",\"photoUrls\": [\"velit mollit dolore\",\"sed\"],\"id\": 1337,\"tags\": [{\"id\": -15522919,\"name\": \"velit ut in esse aliquip\"},{\"id\": 36974757,\"name\": \"ullamco mollit sed commodo\"}],\"status\": \"available\"}"
	h := main.BuildHandler()
	r := httptest.NewRequest("POST", "/api/pet", strings.NewReader(existingreq))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	var tests = []struct {
		name     string
		method   string
		url      string
		headers  map[string]string
		body     string
		wantcode int
		wantbody string
	}{
		{
			"sanity-add",
			"POST",
			"/api/pet",
			map[string]string{"Content-Type": "application/json"},
			"{\"name\": \"mypet\",\"photoUrls\": [\"velit mollit dolore\",\"sed\"],\"id\": 10,\"tags\": [{\"id\": -15522919,\"name\": \"velit ut in esse aliquip\"},{\"id\": 36974757,\"name\": \"ullamco mollit sed commodo\"}],\"status\": \"available\"}",
			200,
			"{\"id\":10, \"name\":\"mypet\", \"photoUrls\":null}",
		},
		{
			"add-pet-empty",
			"POST",
			"/api/pet",
			map[string]string{},
			"{}",
			400,
			"{\"error\": \"Bad Request\"}",
		},
		{
			"already-exists",
			"POST",
			"/api/pet",
			map[string]string{"Content-Type": "application/json"},
			fmt.Sprintf("{\"name\": \"mypet\",\"photoUrls\": [\"velit mollit dolore\",\"sed\"],\"id\": %d,\"tags\": [{\"id\": -15522919,\"name\": \"velit ut in esse aliquip\"},{\"id\": 36974757,\"name\": \"ullamco mollit sed commodo\"}],\"status\": \"available\"}", existingpet.Id),
			409,
			"{\"error\": \"Conflict\"}",
		},
		{
			"sanity-get",
			"GET",
			fmt.Sprintf("/api/pet/%d", existingpet.Id),
			map[string]string{"Content-Type": "application/json"},
			"",
			200,
			fmt.Sprintf("{\"id\":%d, \"name\":\"%s\", \"photoUrls\":null}", existingpet.Id, existingpet.Name),
		},
		{
			"sanity-not-found",
			"GET",
			fmt.Sprintf("/api/pet/%d", existingpet.Id+1),
			map[string]string{"Content-Type": "application/json"},
			"",
			404,
			"{\"error\": \"Not Found\"}",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			for k, v := range tt.headers {
				r.Header.Add(k, v)
			}
			h.ServeHTTP(w, r)
			require.Equal(t, tt.wantcode, w.Code)
			body, _ := io.ReadAll(w.Body)
			require.JSONEq(t, tt.wantbody, string(body))
		})
	}
}

func TestEndToEndUnimplemented(t *testing.T) {
	h := main.BuildHandler()
	var tests = []struct {
		name   string
		method string
		url    string
		body   string
	}{
		{"pet-upload-image", "POST", "/api/pet/1/uploadImage", ""},
		{"pet-update", "PUT", "/api/pet", "{}"},
		{"pet-find-by-status", "GET", "/api/pet/findByStatus", ""},
		{"pet-find-by-tags", "GET", "/api/pet/findByTags", ""},
		{"pet-update-id", "POST", "/api/pet/1", ""},
		{"pet-delete", "DELETE", "/api/pet/1", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			require.Equal(t, 501, w.Code, "you're not expected to implement the service's logic")
		})
	}
}
