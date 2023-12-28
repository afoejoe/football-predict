package main

import (
	"fmt"
	"net/http"

	"github.com/afoejoe/football-predict/assets"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(app.notFound)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handler("GET", "/static/*filepath", fileServer)

	mux.HandlerFunc("GET", "/", app.home)
	mux.HandlerFunc("GET", "/prediction/:slug", app.single)
	fmt.Println("hello")
	//ADMIN
	mux.Handler("GET", "/admin", app.requireBasicAuthentication(http.HandlerFunc(app.admin)))
	mux.Handler("GET", "/admin/prediction/:id", app.requireBasicAuthentication(http.HandlerFunc(app.editOrCreatePrediction)))

	mux.Handler(http.MethodPost, "/admin/prediction", app.requireBasicAuthentication(http.HandlerFunc(app.createPredictionPost)))
	mux.Handler("DELETE", "/admin/prediction/:id", app.requireBasicAuthentication(http.HandlerFunc(app.deletePrediction)))

	return app.logAccess(app.recoverPanic(app.securityHeaders(mux)))
}
