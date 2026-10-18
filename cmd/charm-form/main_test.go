package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	charmform "github.com/junhinhow/charm-form"
)

func TestRunSemArquivo(t *testing.T) {
	var stdout, stderr strings.Builder
	err := run([]string{}, &stdout, &stderr)
	if err == nil {
		t.Fatal("esperado erro quando nenhum arquivo é informado, obtido nil")
	}
	if !strings.Contains(err.Error(), "-file") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestRunArquivoInexistente(t *testing.T) {
	var stdout, stderr strings.Builder
	err := run([]string{"-file", "nao-existe.yaml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("esperado erro para arquivo inexistente, obtido nil")
	}
	if !strings.Contains(err.Error(), "erro ao ler arquivo") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestRunYAMLInvalido(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, "invalido.yaml")
	if err := os.WriteFile(caminho, []byte("title: [quebrado\n"), 0o644); err != nil {
		t.Fatalf("não foi possível criar o arquivo de teste: %v", err)
	}

	var stdout, stderr strings.Builder
	err := run([]string{"-file", caminho}, &stdout, &stderr)
	if err == nil {
		t.Fatal("esperado erro para YAML inválido, obtido nil")
	}
	if !strings.Contains(err.Error(), "erro ao parsear YAML") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestRunAjuda(t *testing.T) {
	var stdout, stderr strings.Builder
	err := run([]string{"-h"}, &stdout, &stderr)
	if err == nil {
		t.Skip("a flag -h encerrou com sucesso")
	}
	if !strings.Contains(stderr.String(), "charm-form") {
		t.Errorf("a ajuda deveria ser exibida, obtido:\n%s", stderr.String())
	}
}

func TestEscreveResumo(t *testing.T) {
	spec := charmform.FormSpec{
		Title: "T",
		Fields: []charmform.FieldSpec{
			{Name: "nome", Type: charmform.TypeInput},
			{Name: "ativo", Type: charmform.TypeConfirm},
			{Name: "tags", Type: charmform.TypeMultiSelect},
		},
	}
	respostas := charmform.Answers{
		"nome":  "Ana",
		"ativo": true,
		"tags":  []string{"a", "b"},
	}

	var b strings.Builder
	if err := escreveResumo(&b, spec, respostas); err != nil {
		t.Fatalf("escreveResumo retornou erro inesperado: %v", err)
	}
	saida := b.String()
	for _, esperado := range []string{"nome: Ana", "ativo: sim", "tags: a, b"} {
		if !strings.Contains(saida, esperado) {
			t.Errorf("o resumo deveria conter %q, obtido:\n%s", esperado, saida)
		}
	}
}

func TestFormatarValor(t *testing.T) {
	casos := []struct {
		valor    any
		esperado string
	}{
		{nil, ""},
		{"texto", "texto"},
		{true, "sim"},
		{false, "não"},
		{[]string{"a", "b"}, "a, b"},
		{42, "42"},
	}
	for _, c := range casos {
		obtido := formatarValor(c.valor)
		if obtido != c.esperado {
			t.Errorf("formatarValor(%#v) = %q, esperado %q", c.valor, obtido, c.esperado)
		}
	}
}
