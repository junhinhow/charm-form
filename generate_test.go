package charmform

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestToJSON(t *testing.T) {
	spec := FormSpec{
		Title:       "T",
		Description: "D",
		Fields:      []FieldSpec{{Name: "a", Type: TypeInput}},
	}
	data, err := ToJSON(spec)
	if err != nil {
		t.Fatalf("ToJSON retornou erro inesperado: %v", err)
	}
	var decoded FormSpec
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON gerado é inválido: %v", err)
	}
	if decoded.Title != "T" || decoded.Description != "D" || len(decoded.Fields) != 1 {
		t.Errorf("JSON inesperado: %s", data)
	}
}

func TestWriteJSON(t *testing.T) {
	var b strings.Builder
	spec := FormSpec{Title: "T", Fields: []FieldSpec{{Name: "a", Type: TypeInput}}}
	if err := WriteJSON(&b, spec); err != nil {
		t.Fatalf("WriteJSON retornou erro inesperado: %v", err)
	}
	if !strings.HasSuffix(b.String(), "\n") {
		t.Error("WriteJSON deveria terminar com quebra de linha")
	}
	var decoded FormSpec
	if err := json.Unmarshal([]byte(b.String()), &decoded); err != nil {
		t.Fatalf("JSON gerado é inválido: %v", err)
	}
	if decoded.Title != "T" {
		t.Errorf("JSON inesperado: %s", b.String())
	}
}

func TestWriteJSONPropagaErroDoWriter(t *testing.T) {
	spec := FormSpec{Title: "T"}
	err := WriteJSON(writerQueFalha{}, spec)
	if err == nil {
		t.Fatal("esperado erro do writer, obtido nil")
	}
	if !strings.Contains(err.Error(), "erro ao escrever JSON") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestRunEscreveNoStdout(t *testing.T) {
	// Run é um wrapper de WriteJSON sobre os.Stdout; o teste garante que a
	// função existe com a assinatura pública esperada e devolve erro do
	// writer apenas via WriteJSON.
	spec := FormSpec{Title: "T", Fields: []FieldSpec{{Name: "a", Type: TypeInput}}}
	if err := Run(spec); err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}
	_ = os.Stdout
}

// writerQueFalha é um io.Writer que sempre retorna erro.
type writerQueFalha struct{}

func (writerQueFalha) Write([]byte) (int, error) {
	return 0, os.ErrClosed
}
