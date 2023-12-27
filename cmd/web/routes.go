package main

import (
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
	mux.Handler("GET", "/admin", app.requireBasicAuthentication(http.HandlerFunc(app.admin)))
	mux.Handler("GET", "/admin/:slug", app.requireBasicAuthentication(http.HandlerFunc(app.editOrCreatePrediction)))
	// mux.HandlerFunc("GET", "/admin", app.requireBasicAuthentication(http.HandlerFunc(app.admin)))

	mux.Handler(http.MethodPost, "/admin/create", app.requireBasicAuthentication(http.HandlerFunc(app.createPredictionPost)))
	mux.HandlerFunc("GET", "/prediction/:slug", app.single)

	return app.logAccess(app.recoverPanic(app.securityHeaders(mux)))
}
