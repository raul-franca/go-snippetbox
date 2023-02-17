package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	//Cria nova instância da nossa estrutura de aplicação. Por enquanto, isso só
	//contém apenas alguns mock loggers (which discard anything written to them).
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
	// Cria um httptest.NewTLSServer() para servidor,
	// passando o valor retornado pelo nosso método app.routes() como o
	// manipulador para o servidor. Isso inicia um servidor HTTPS que escuta em uma
	// porta escolhida aleatoriamente durante o teste.
	// Observe que adiamos uma chamada para ts.Close() para desligar o servidor quando
	// o teste termina.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// O endereço de rede que o servidor de teste está escutando está contido
	// no campo ts.URL. Podemos usar isso junto com o ts.Client().Get()
	// método para fazer uma solicitação GET /ping no servidor de teste. Esse
	// retorna uma estrutura http.Response contendo a resposta.
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	// Podemos então verificar o valor do código de status da resposta e do corpo usando
	//o mesmo código de antes.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
