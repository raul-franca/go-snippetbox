package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extrai o valor da query string e tenta convert para integer com strconv.Atoi()
	// convert it to an integer using the strconv.Atoi() function.
	// se nao convert ou o valor for menor que 1,retorna 404 page // not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	//and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// content as the response body. The last parameter to Execute() represents any
	//dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // avisa que o method POST eeh permitido
		w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte("Method não autorizado"))
		http.Error(w, "Method não autorizado", 405)
		return
	}

	_, err := w.Write([]byte("Create a new snippet..."))
	if err != nil {
		return
	}
}
