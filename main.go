package main

import (
	"com.gitlab.pscheid92/idcardgenerator/pkg/models"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func process(w http.ResponseWriter, r *http.Request) {
	model := models.NewViewModel()

	// use x-forwarded-prefix header if set
	if pathPrefix := r.Header.Get("X-Forwarded-Prefix"); pathPrefix != "" {
		model.PathPrefix = pathPrefix
	}

	switch r.Method {
	case "GET":
		model.CardOptions[0].Selected = true
		_ = templates.ExecuteTemplate(w, "index.html", model)
	case "POST":
		// get values
		model.Birthday = r.FormValue("birthday")
		model.Expiration = r.FormValue("expiration")
		model.Manipulation = r.FormValue("manipulation") == "manipulation"

		// recreate card type and calculate perso data
		switch r.FormValue("cardtype") {
		case "newid":
			model.CardOptions[0].Selected = true
			model.CalculateNewId()
		case "oldid":
			model.CardOptions[1].Selected = true
			model.CalculateOldId()
		case "passport":
			model.CardOptions[2].Selected = true
			model.CalculatePassport()
		}

		_ = templates.ExecuteTemplate(w, "index.html", model)
	}
}

func main() {
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", process)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
