package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	forms "github.com/hbourgeot/portfolio/internal/models"
)

var NameForm, EmailForm, SubjectForm, MessageForm string

func Redirect(w http.ResponseWriter, r *http.Request, location string) {
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func home(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.URL.Path != "/" {
		Redirect(w, r, "/not-found")
		return
	}
	temp := template.Must(template.ParseFiles("go/src/portfolio/ui/index.html"))
	err := temp.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("go/src/portfolio/ui/404.html"))
	err := temp.Execute(w, nil)
	if err != nil {
		return
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func SendForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	NameForm, EmailForm, SubjectForm, MessageForm = r.PostForm.Get("name-form"), r.PostForm.Get("email-form"), r.PostForm.Get("subject-form"), r.PostForm.Get("message-form")
	defer Redirect(w, r, "/")
	fmt.Fprintf(w, "The form values has been send. I will contact you soon.")
	if err != nil {
		return
	}
	err = forms.Insert(NameForm, EmailForm, SubjectForm, MessageForm)
	if err != nil {
		log.Printf("algo hiciste mal, %s", err)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("go/src/portfolio/ui/login.html"))
	err := temp.Execute(w, nil)
	if err != nil {
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}
	username, pass := r.PostForm.Get("username-login"), r.PostForm.Get("password-login")
	res, err := forms.CheckLogin(username, pass)
	if err != nil {
		log.Fatal(err)
		return
	}
	confirmLogin(w, r, res)
}

func showMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := forms.ShowMSGs()
	if err != nil {
		log.Fatal(err)
		return
	}
	enableCors(&w)
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		return
	}
}

func confirmLogin(w http.ResponseWriter, r *http.Request, res bool) {
	if res != false {
		http.HandleFunc("/json", showMessages)
		Redirect(w, r, "/admin/panel")
	}
}

func ShowPanel(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("go/src/ṕortfolio/ui/panel.html"))
	err := temp.Execute(w, nil)
	if err != nil {
		return
	}
}
