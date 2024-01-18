package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/afoejoe/football-predict/internal/database"
	"github.com/afoejoe/football-predict/internal/funcs"
	"github.com/afoejoe/football-predict/internal/request"
	"github.com/afoejoe/football-predict/internal/response"
	"github.com/afoejoe/football-predict/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type PredictionForm struct {
	ID             int64               `form:"id"`
	Title          string              `form:"title"`
	Body           string              `form:"body"`
	Keywords       string              `form:"keywords"`
	ScheduledAt    time.Time           `form:"scheduled_at"`
	Odds           float64             `form:"odds"`
	PredictionType string              `form:"prediction_type"`
	IsFeatured     bool                `form:"is_featured"`
	IsArchived     bool                `form:"is_archived"`
	LeagueID       int64               `form:"league_id"`
	League         database.League     `form:"-"`
	Validator      validator.Validator `form:"-"`
}

var ()

func (app *application) editOrCreatePrediction(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("id")

	if slug == "" {
		app.badRequest(w, r, errors.New("no slug provided"))
		return
	}

	leagues, err := app.db.GetLeagues()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data["Leagues"] = leagues

	if slug == "new" {
		data["Form"] = PredictionForm{
			Title:          "",
			PredictionType: "",
			ScheduledAt:    time.Now(),
		}
		err := response.Page(w, http.StatusOK, data, PredictionCreateTemplate)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	id, err := strconv.Atoi(slug)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	prediction, err := app.db.GetPrediction(int64(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Form"] = PredictionForm{
		ID:             prediction.ID,
		Title:          prediction.Title,
		Body:           prediction.Body,
		Keywords:       prediction.Keywords,
		ScheduledAt:    prediction.ScheduledAt,
		Odds:           prediction.Odds,
		PredictionType: prediction.PredictionType,
		IsFeatured:     prediction.IsFeatured,
		IsArchived:     prediction.IsArchived,
		League:         prediction.League,
		LeagueID:       prediction.LeagueID,
	}
	//convert prediction to PredictionForm

	err = response.Page(w, http.StatusOK, data, PredictionCreateTemplate)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createPredictionPost(w http.ResponseWriter, r *http.Request) {
	form := PredictionForm{}
	err := request.DecodePostForm(r, &form)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	form.Validator.CheckField(form.Title != "", "title", "must be provided")
	form.Validator.CheckField(form.ScheduledAt != time.Time{}, "scheduled_at", "must be provided")
	if form.Validator.HasErrors() {
		data := app.newTemplateData(r)
		data["Form"] = form
		err = response.Page(w, http.StatusUnprocessableEntity, data, PredictionCreateTemplate)

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
		LeagueID:       form.LeagueID,
		Slug:           funcs.Slugify(form.Title + " " + form.PredictionType + " " + form.ScheduledAt.Format("2006-01-02")),
	}
	if form.ID != 0 {
		p.ID = form.ID
		err = app.db.UpdatePrediction(p)
	} else {
		err = app.db.InsertPrediction(p)
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) deletePrediction(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	err = app.db.DeletePrediction(int64(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// func (app *application) protected(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("This is a protected handler"))
// }
