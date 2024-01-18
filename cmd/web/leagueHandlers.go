package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/afoejoe/football-predict/internal/database"
	"github.com/afoejoe/football-predict/internal/request"
	"github.com/afoejoe/football-predict/internal/response"
	"github.com/afoejoe/football-predict/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type LeagueForm struct {
	ID        int64               `form:"id"`
	Title     string              `form:"title"`
	Validator validator.Validator `form:"-"`
}

func (app *application) adminLeague(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	leagues, err := app.db.GetLeagues()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Leagues"] = leagues
	err = response.Page(w, http.StatusOK, data, "pages/admin-league.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) editOrCreateLeague(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("id")

	if slug == "" {
		app.badRequest(w, r, errors.New("no slug provided"))
		return
	}

	if slug == "new" {
		data["Form"] = LeagueForm{
			Title: "",
		}
		err := response.Page(w, http.StatusOK, data, LeagueCreateTemplate)
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

	prediction, err := app.db.GetLeague(int64(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Form"] = LeagueForm{
		ID:    prediction.ID,
		Title: prediction.Title,
	}

	err = response.Page(w, http.StatusOK, data, LeagueCreateTemplate)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createLeaguePost(w http.ResponseWriter, r *http.Request) {
	form := LeagueForm{}
	err := request.DecodePostForm(r, &form)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	form.Validator.CheckField(form.Title != "", "title", "must be provided")

	if form.Validator.HasErrors() {
		data := app.newTemplateData(r)
		data["Form"] = form
		err = response.Page(w, http.StatusUnprocessableEntity, data, LeagueCreateTemplate)

		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	p := &database.League{
		Title: form.Title,
	}
	if form.ID != 0 {
		p.ID = form.ID
		err = app.db.UpdateLeague(p)
	} else {
		err = app.db.InsertLeague(p)
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/league", http.StatusSeeOther)
}

func (app *application) deleteLeague(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	err = app.db.DeleteLeague(int64(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/admin/league", http.StatusSeeOther)
}

// func (app *application) protected(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("This is a protected handler"))
// }
