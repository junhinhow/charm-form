package renderer

import (
	"io"
	"strings"
	"testing"

	charmform "github.com/junhinhow/charm-form"
)

// TestRunComEntradaProgramada executa o formulário completo com uma entrada
// programada, verificando que as respostas coletadas chegam ao chamador.
func TestRunComEntradaProgramada(t *testing.T) {
	spec := charmform.FormSpec{
		Title: "Formulário",
		Fields: []charmform.FieldSpec{
			{Name: "nome", Type: charmform.TypeInput, Required: true},
			{Name: "idade", Type: charmform.TypeInput, Validate: "^[0-9]+$"},
		},
	}

	// "Ana\r" preenche o primeiro campo e confirma; "30\r" preenche o segundo.
	entrada := strings.NewReader("Ana\r30\r")

	respostas, err := RunWithOptions(spec, Options{
		Input:  entrada,
		Output: io.Discard,
		Width:  80,
		Height: 24,
	})
	if err != nil {
		t.Fatalf("RunWithOptions retornou erro inesperado: %v", err)
	}

	if v, ok := respostas["nome"].(string); !ok || v != "Ana" {
		t.Errorf("resposta 'nome' inesperada: %#v", respostas["nome"])
	}
	if v, ok := respostas["idade"].(string); !ok || v != "30" {
		t.Errorf("resposta 'idade' inesperada: %#v", respostas["idade"])
	}
}

// TestRunCancela verifica que o cancelamento é reportado ao chamador.
func TestRunCancela(t *testing.T) {
	spec := charmform.FormSpec{
		Title:  "Formulário",
		Fields: []charmform.FieldSpec{{Name: "nome", Type: charmform.TypeInput}},
	}

	// "\x03" é ctrl+c.
	_, err := RunWithOptions(spec, Options{
		Input:  strings.NewReader("\x03"),
		Output: io.Discard,
		Width:  80,
		Height: 24,
	})
	if err == nil {
		t.Fatal("esperado erro de cancelamento, obtido nil")
	}
	if err != ErrCancelado {
		t.Errorf("erro esperado %v, obtido %v", ErrCancelado, err)
	}
}
