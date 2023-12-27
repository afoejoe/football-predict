package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/afoejoe/football-predict/internal/database"
	"github.com/afoejoe/football-predict/internal/funcs"
	"github.com/afoejoe/football-predict/internal/request"
	"github.com/afoejoe/football-predict/internal/response"
	"github.com/afoejoe/football-predict/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type PredictionForm struct {
	Title          string              `form:"title"`
	Body           string              `form:"body"`
	Keywords       string              `form:"keywords"`
	ScheduledAt    time.Time           `form:"scheduled_at"`
	Odds           float64             `form:"odds"`
	PredictionType string              `form:"prediction_type"`
	IsFeatured     bool                `form:"is_featured"`
	IsArchived     bool                `form:"is_archived"`
	Validator      validator.Validator `form:"-"`
}

func (app *application) editOrCreatePrediction(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("slug")

	if slug == "" {
		app.badRequest(w, r, errors.New("no slug provided"))
		return
	}

	if slug == "new" {
		data["Form"] = PredictionForm{
			Title: "",
		}
		err := response.Page(w, http.StatusOK, data, "pages/admin-create.html")
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	prediction, err := app.db.GetPredictionBySlug(slug)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Form"] = prediction
	err = response.Page(w, http.StatusOK, data, "pages/admin-edit.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createPredictionPost(w http.ResponseWriter, r *http.Request) {
	// data := app.newTemplateData(r)
	// data["Form"] = database.Prediction{}

	form := PredictionForm{}
	err := request.DecodePostForm(r, &form)

	if err != nil {
		fmt.Println(err)
		app.badRequest(w, r, err)
		// set header here
		w.Header().Set("HX-Retarget", "error")
		return
	}

	form.Validator.CheckField(form.Title != "", "title", "must be provided")
	form.Validator.CheckField(form.ScheduledAt != time.Time{}, "scheduled_at", "must be provided")
	if form.Validator.HasErrors() {
		data := app.newTemplateData(r)
		data["Form"] = form
		err = response.Page(w, http.StatusUnprocessableEntity, data, "pages/admin-create.html")

		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	p := &database.Prediction{
		Title:          form.Title,
		Body:           form.Body,
		Keywords:       form.Keywords,
		ScheduledAt:    form.ScheduledAt,
		Odds:           form.Odds,
		PredictionType: form.PredictionType,
		IsFeatured:     form.IsFeatured,
		IsArchived:     form.IsArchived,
		Slug:           funcs.Slugify(form.Title),
	}

	fmt.Println(p)
	err = app.db.InsertPrediction(p)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// func (app *application) protected(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("This is a protected handler"))
// }
