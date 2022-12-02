package main

import (
	"errors"
	"fmt"
	"github.com/raul-franca/go-snippetbox/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w) //notFound() helper
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	//
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	//// template.ParseFiles() le os template.
	////the http.Error() function to send a generic 500 Internal Server Error
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err) // serverError() helper.
	//	return
	//}
	//
	//err = ts.Execute(w, nil)
	//if err != nil {
	//	app.serverError(w, err) // Use the serverError() helper.
	//}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extrai o valor da query string e tenta convert para integer com strconv.Atoi()
	// converta-o em um número inteiro usando a função strconv.Atoi().
	// Se nao convert ou o valor for menor que 1,retorna 404 page // not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	//specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Escreva o snippet como plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost) // avisa que o method POST eeh permitido
	//	w.WriteHeader(http.StatusMethodNotAllowed)
	//	//w.Write([]byte("Method não autorizado"))
	//	//http.Error(w, "Method não autorizado", 405)
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	return
	//}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"
	// Pass the data to the SnippetModel.Insert() method, receiving the // ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
