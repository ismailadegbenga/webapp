package main

import (
	"github.com/ismailadegbenga/webapp/models"

	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	tmpl = template.New("")
)

func init() {
	// templating

    tmpl.Funcs(template.FuncMap{"StringsJoin": strings.Join})
    _, err = tmpl.ParseGlob(filepath.Join(".", "templates", "*.html"))
    if err != nil {
        log.Fatalf("Unable to parse templates: %v\n", err)
    }
}


func GetContacts(w http.ResponseWriter, r *http.Request) {
	var cs models.Contacts
	err := models.cs.Get()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
	}
	tmpl.ExecuteTemplate(w, "index.html", cs)
}

func main() {
	var c models.Contacts{}
	http.HandlerFunc("/", GetContacts())
}
