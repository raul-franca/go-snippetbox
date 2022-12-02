package main

import (
	"errors"
	"fmt"
	"github.com/raul-franca/go-snippetbox/pkg/models"
	"html/template"
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

	// Create an instance of a templateData struct holding the slice de snippets.
	data := &templateData{Snippets: s}
	files := []string{"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Passe a estrutura templateData ao executar o modelo.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// strconv.Atoi() converte para um int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	//Buscar o registro pelo id
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//Crie uma instância de uma templateData contendo os dados do snippet.
	data := &templateData{Snippet: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

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
	// Passa os dados para o método SnippetModel.Insert(), recebendo o
	//ID do novo registro de volta.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirecione o usuário para a página relevante do snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
