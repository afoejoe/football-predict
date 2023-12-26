package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/afoejoe/football-predict/internal/database"
	"github.com/afoejoe/football-predict/internal/response"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	predictions, err := app.db.GetPredictions()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	featured, err := app.db.GetFeatured()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			featured = []*database.Prediction{}
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	// data["Predictions"] = []*database.Prediction{}
	data["Predictions"] = predictions
	data["Featured"] = featured

	err = response.Page(w, http.StatusOK, data, "pages/home.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) single(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("slug")
	
	if slug == "" {
		app.serverError(w, r, errors.New("No slug provided"))
		return
	}

	prediction, err := app.db.GetPredictionBySlug(slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data["Prediction"] = prediction

	err = response.Page(w, http.StatusOK, data, "pages/single.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) admin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/admin-home.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

// func (app *application) protected(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("This is a protected handler"))
// }
