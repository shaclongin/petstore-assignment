package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tanninio/home-assignment/internal/adapters"
	"github.com/tanninio/home-assignment/internal/app"
	"github.com/tanninio/home-assignment/internal/common"
	ports "github.com/tanninio/home-assignment/internal/ports/http"
)

func ConfigureRouters(root, svcrouter *mux.Router) {
	svcrouter.Use(common.LoggingMiddleware())
}

func BuildHandler() http.Handler {
	app := app.NewApplication(adapters.NewMemRepository())
	return ports.HttpCreateServiceHandler(app, "/api", ConfigureRouters)
}

func main() {
	ports.HttpServeHandler(":8080", BuildHandler())
}
