package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"henrry.online/internal/validator"
)

type messageForm struct {
	Name                string `form:"name-form"`
	Email               string `form:"email-form"`
	Subject             string `form:"subject-form"`
	Message             string `form:"message-form"`
	validator.Validator `form:"-"`
}

type loginForm struct {
	Username            string `form:"username-login"`
	Password            string `form:"password-login"`
	validator.Validator `form:"-"`
}

func (folio *portfolio) home(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if strings.Contains(r.URL.Path, "/api/create/hangman") || strings.Contains(r.URL.Path, "/api/get/hangman") {
		return
	}

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/not-found", http.StatusNotFound)
		return
	}

	data := folio.newTemplateData(r)
	folio.render(w, http.StatusOK, "./ui/index.gohtml", data)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (folio *portfolio) SendMessage(w http.ResponseWriter, r *http.Request) {
	var form messageForm

	err := folio.decodePostForm(r, &form)
	if err != nil {
		folio.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Subject), "subject", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Message), "message", "This field cannot be blank")

	form.CheckField(validator.MinChars(form.Name, 5), "name", "This field must be 5 characters long")
	form.CheckField(validator.MinChars(form.Email, 20), "email", "This field must be 20 characters long")
	form.CheckField(validator.MinChars(form.Subject, 50), "subject", "This field must be 50 characters long")
	form.CheckField(validator.MaxChars(form.Message, 200), "message", "This field cannot be more than 150 characters long")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please, enter a valid email")

	if !form.Valid() {
		data := folio.newTemplateData(r)
		data.Form = form
		folio.render(w, http.StatusUnprocessableEntity, "./ui/index.gohtml", data)
		return
	}

	err = folio.messages.Insert(form.Name, form.Email, form.Subject, form.Message)
	if err != nil {
		folio.serverError(w, err)
		return
	}

	folio.sessionManager.Put(r.Context(), "flash", "Form successfully submited!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (folio *portfolio) Login(w http.ResponseWriter, r *http.Request) {
	data := folio.newTemplateData(r)
	data.Form = nil
	folio.render(w, http.StatusOK, "./ui/login.gohtml", data)
}

func (folio *portfolio) LoginPost(w http.ResponseWriter, r *http.Request) {
	var form loginForm

	err := folio.decodePostForm(r, &form)
	if err != nil {
		folio.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Username, 20), "username", "This field must be 20 characters long")
	form.CheckField(validator.MinChars(form.Password, 20), "password", "This field must be 20 characters long")
	form.CheckField(validator.Matches(form.Username, validator.EmailRX), "username", "Please, enter a valid email")

	if !form.Valid() {
		data := folio.newTemplateData(r)
		data.Form = form
		folio.render(w, http.StatusUnprocessableEntity, "./ui/login.gohtml", data)
		return
	}

	err = folio.users.CheckLogin(form.Username, form.Password)
	if err != nil {
		data := folio.newTemplateData(r)
		data.Form = form
		folio.sessionManager.Put(r.Context(), "flash", "User or password may be incorrect")
		folio.render(w, http.StatusUnprocessableEntity, "./ui/login.gohtml", data)
		return
	}

	http.Redirect(w, r, "/panel/yes", http.StatusSeeOther)
}

func (folio *portfolio) ShowPanel(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if logged := params.ByName("logged"); logged != "yes" {
		folio.notFound(w)
		return
	}

	messages, err := folio.messages.ShowMSGs()
	if err != nil {
		folio.serverError(w, err)
		return
	}

	data := folio.newTemplateData(r)
	data.Messages = messages

	folio.render(w, http.StatusOK, "./ui/panel.gohtml", data)
}

func (folio *portfolio) StoreHome(w http.ResponseWriter, r *http.Request) {
	folio.render(w, http.StatusOK, "./ui/store.gohtml", nil)
}

func (folio *portfolio) StoreLogin(w http.ResponseWriter, r *http.Request) {
	folio.render(w, http.StatusOK, "./ui/store-login.gohtml", nil)
}
