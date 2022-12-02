package main

import "github.com/raul-franca/go-snippetbox/pkg/models"

// Definir um tipo templateData para atuar como a estrutura de retenção para
// quaisquer dados dinâmicos que queremos passar para nossos modelos HTML.
// No momento contém apenas um campo, mas adicionaremos mais
// para ele à medida que a construção progride.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
