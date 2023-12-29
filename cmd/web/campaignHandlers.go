package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/afoejoe/football-predict/internal/request"
	"github.com/afoejoe/football-predict/internal/validator"
	"github.com/getbrevo/brevo-go/lib"
	"github.com/julienschmidt/httprouter"
)

type subscriptionForm struct {
	Email     string              `form:"email"`
	Validator validator.Validator `form:"-"`
}

const (
	ListId = 9
)

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	form := subscriptionForm{}
	err := request.DecodeForm(r, &form)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	validator.ValidateEmail(form.Validator, form.Email)

	if form.Validator.HasErrors() {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, _, err = app.brevoClient.ContactsApi.CreateContact(ctx, lib.CreateContact{
		ListIds:       []int64{ListId},
		Email:         form.Email,
		UpdateEnabled: true,
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	_, err = w.Write([]byte("Thank you for subscribing"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) sendCampaign(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	fmt.Println(err, id)

	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, _, err = app.brevoClient.EmailCampaignsApi.GetEmailCampaign(ctx, int64(4), &lib.GetEmailCampaignOpts{})
	// fmt.Println(c)
	s := &lib.CreateEmailCampaignSender{
		Name:  "MyName",
		Email: "newsletter@naijarank.com",
	}

	// html string content from template emails/example.html file

	prediction, err := app.db.GetPrediction(int64(id))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	filePath := "assets/emails/prediction.html"
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		fmt.Println(err, "Error3")

		app.serverError(w, r, err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, prediction); err != nil {
		fmt.Println(err, "Error2")
		app.serverError(w, r, err)
		return
	}

	fmt.Println("content", tpl.String())
	params3 := lib.CreateEmailCampaign{
		Sender:                s,
		InlineImageActivation: false,
		Name:                  "Prediction " + prediction.Title + " " + time.Now().String(),
		SendAtBestTime:        false,
		AbTesting:             false,
		IpWarmupEnable:        false,
		// TemplateId:            int64(3),
		HtmlContent: tpl.String(),
		Subject:     "New Prediction Just Now!",
		Recipients: &lib.CreateEmailCampaignRecipients{
			ListIds: []int64{9},
		},
	}

	a, _, err := app.brevoClient.EmailCampaignsApi.CreateEmailCampaign(ctx, params3)

	if err != nil {
		fmt.Println(err, "Error")
		app.serverError(w, r, err)
		return
	}

	fmt.Println(err, a, "Sent")
	_, err = app.brevoClient.EmailCampaignsApi.SendEmailCampaignNow(ctx, a.Id)
	fmt.Println(err, a, "Sent 2")

	if err != nil {

		app.serverError(w, r, err)
		return

	}

	prediction.Campaigned = true
	err = app.db.UpdatePrediction(prediction)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
