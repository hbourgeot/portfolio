package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/hbourgeot/portfolio/internal/forms"
)

type templateData struct {
	Messages []*forms.Message
	Form     any
	Flash    string
}

func (folio *portfolio) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		Flash: folio.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (folio *portfolio) checkTemplate(w http.ResponseWriter, dir string) *template.Template {
	temp := template.Must(template.ParseFiles(dir))
	return temp
}

func (folio *portfolio) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts := folio.checkTemplate(w, page)

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, data)
	if err != nil {
		folio.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	buf.WriteTo(w)
}
