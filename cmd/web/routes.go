package main

import (
	"net/http"

	"github.com/afoejoe/football-predict/assets"

	"github.com/julienschmidt/httprouter"
)

var (
	PredictionCreateTemplate = "pages/admin-prediction-create.html"
	LeagueCreateTemplate     = "pages/admin-league-create.html"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(app.notFound)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handler("GET", "/static/*filepath", fileServer)

	mux.HandlerFunc("GET", "/", app.home)
	mux.HandlerFunc("GET", "/prediction/:slug", app.single)

	// CAMPAIGN
	mux.HandlerFunc("POST", "/subscribe", app.subscribe)
	mux.HandlerFunc("POST", "/campaign/:id", app.sendCampaign)

	/**
	 * ADMIN
	 */
	// Admin Home
	mux.Handler("GET", "/admin", app.requireBasicAuthentication(http.HandlerFunc(app.admin)))

	// Admin Prediction
	mux.Handler("GET", "/admin/prediction/:id", app.requireBasicAuthentication(http.HandlerFunc(app.editOrCreatePrediction)))
	mux.Handler(http.MethodPost, "/admin/prediction", app.requireBasicAuthentication(http.HandlerFunc(app.createPredictionPost)))
	mux.Handler("DELETE", "/admin/prediction/:id", app.requireBasicAuthentication(http.HandlerFunc(app.deletePrediction)))

	// Admin League
	mux.Handler("GET", "/admin/league", app.requireBasicAuthentication(http.HandlerFunc(app.adminLeague)))
	mux.Handler("GET", "/admin/league/:id", app.requireBasicAuthentication(http.HandlerFunc(app.editOrCreateLeague)))
	mux.Handler(http.MethodPost, "/admin/league", app.requireBasicAuthentication(http.HandlerFunc(app.createLeaguePost)))
	mux.Handler("DELETE", "/admin/league/:id", app.requireBasicAuthentication(http.HandlerFunc(app.deleteLeague)))

	return app.logAccess(app.recoverPanic(app.securityHeaders(mux)))
}
