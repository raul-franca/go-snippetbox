package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w) //notFound() helper
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// template.ParseFiles() le os template.
	//the http.Error() function to send a generic 500 Internal Server Error
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // serverError() helper.
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
	}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extrai o valor da query string e tenta convert para integer com strconv.Atoi()
	// convert it to an integer using the strconv.Atoi() function.
	// se nao convert ou o valor for menor que 1,retorna 404 page // not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // avisa que o method POST eeh permitido
		w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte("Method não autorizado"))
		//http.Error(w, "Method não autorizado", 405)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := w.Write([]byte("Create a new snippet..."))
	if err != nil {
		return
	}
}
