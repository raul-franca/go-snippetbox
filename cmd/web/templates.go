package main

import (
	"github.com/raul-franca/go-snippetbox/pkg/forms"
	"github.com/raul-franca/go-snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Definir um tipo templateData para atuar como a estrutura de retenção para
// quaisquer dados dinâmicos que queremos passar para nossos modelos HTML.
// No momento contém apenas um campo, mas adicionaremos mais
// para ele à medida que a construção progride.
type templateData struct {
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Inicializa um novo mapa para atuar como cache.
	cache := map[string]*template.Template{}

	//O filepath.Glob função para obter um slice de todos o filepaths com
	//// a extensão '.page.tmpl'. Isso basicamente nos dá uma fatia de todos os
	////modelos de 'página' para o aplicativo.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extrai o nome do arquivo (como 'home.page.tmpl') do caminho completo do arquivo
		//e atribua-o à variável name.
		name := filepath.Base(page)
		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use o método ParseGlob para adicionar qualquer modelo de 'layout' ao
		//conjunto de modelo (no nosso caso, é apenas o layout 'base' no momento).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Adiciona o conjunto de templates ao cache, usando o nome da página
		//(como 'home.page.tmpl') como a chave.
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}
