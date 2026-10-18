package charmform

import (
	"path/filepath"
	"strings"
	"testing"
)

// yamlValido é uma especificação de formulário completa e válida.
const yamlValido = `
title: Formulário de Deploy
description: Configurações de publicação
fields:
  - name: environment
    type: select
    label: Ambiente
    options: [staging, production]
    required: true
  - name: version
    type: input
    label: Versão
    placeholder: v1.0.0
    validate: "^v\\d+\\.\\d+\\.\\d+$"
`

func TestLoadFromBytesValido(t *testing.T) {
	spec, err := LoadFromBytes([]byte(yamlValido))
	if err != nil {
		t.Fatalf("LoadFromBytes retornou erro inesperado: %v", err)
	}
	if spec.Title != "Formulário de Deploy" {
		t.Errorf("título esperado %q, obtido %q", "Formulário de Deploy", spec.Title)
	}
	if spec.Description != "Configurações de publicação" {
		t.Errorf("descrição esperada %q, obtida %q", "Configurações de publicação", spec.Description)
	}
	if len(spec.Fields) != 2 {
		t.Fatalf("esperados 2 campos, obtidos %d", len(spec.Fields))
	}
	if spec.Fields[0].Type != TypeSelect {
		t.Errorf("tipo do primeiro campo esperado %q, obtido %q", TypeSelect, spec.Fields[0].Type)
	}
	if !spec.Fields[0].Required {
		t.Error("o primeiro campo deveria ser obrigatório")
	}
	if spec.Fields[1].Validate == "" {
		t.Error("o regex de validação do segundo campo não foi carregado")
	}
}

func TestLoadFromBytesYAMLInvalido(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: [quebrado\nfields:"))
	if err == nil {
		t.Fatal("esperado erro para YAML inválido, obtido nil")
	}
	if !strings.Contains(err.Error(), "erro ao parsear YAML") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesSemTitulo(t *testing.T) {
	_, err := LoadFromBytes([]byte("fields:\n  - name: a\n    type: input\n"))
	if err == nil {
		t.Fatal("esperado erro para formulário sem título, obtido nil")
	}
	if !strings.Contains(err.Error(), "title") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesSemCampos(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: Vazio\nfields: []\n"))
	if err == nil {
		t.Fatal("esperado erro para formulário sem campos, obtido nil")
	}
	if !strings.Contains(err.Error(), "pelo menos um campo") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesCampoSemNome(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - type: input\n"))
	if err == nil {
		t.Fatal("esperado erro para campo sem nome, obtido nil")
	}
	// O índice do campo é 1-based: o usuário conta os campos a partir de 1.
	if !strings.Contains(err.Error(), "campo 1") {
		t.Errorf("o erro deveria citar o campo 1 (base 1), obtido: %v", err)
	}
	if strings.Contains(err.Error(), "campo 0") {
		t.Errorf("o erro não deveria usar índice 0-based, obtido: %v", err)
	}
}

func TestLoadFromBytesCampoSemTipo(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n"))
	if err == nil {
		t.Fatal("esperado erro para campo sem tipo, obtido nil")
	}
	if !strings.Contains(err.Error(), "campo 'a'") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesTipoInvalido(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n    type: radio\n"))
	if err == nil {
		t.Fatal("esperado erro para tipo inválido, obtido nil")
	}
	if !strings.Contains(err.Error(), "inválido") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesTipoNormalizado(t *testing.T) {
	spec, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n    type: Input\n  - name: b\n    type: MULTISELECT\n    options: [x]\n"))
	if err != nil {
		t.Fatalf("LoadFromBytes retornou erro inesperado: %v", err)
	}
	if spec.Fields[0].Type != TypeInput {
		t.Errorf("tipo esperado %q, obtido %q", TypeInput, spec.Fields[0].Type)
	}
	if spec.Fields[1].Type != TypeMultiSelect {
		t.Errorf("tipo esperado %q, obtido %q", TypeMultiSelect, spec.Fields[1].Type)
	}
}

func TestLoadFromBytesSelectSemOpcoes(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n    type: select\n"))
	if err == nil {
		t.Fatal("esperado erro para select sem opções, obtido nil")
	}
	if !strings.Contains(err.Error(), "opção") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromBytesNomesDuplicados(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n    type: input\n  - name: a\n    type: input\n"))
	if err == nil {
		t.Fatal("esperado erro para nomes duplicados, obtido nil")
	}
	if !strings.Contains(err.Error(), "duplicado") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
	// O erro deve apontar o campo duplicado (posição 2, base 1).
	if !strings.Contains(err.Error(), "campo 2") {
		t.Errorf("o erro deveria citar o campo 2 (base 1), obtido: %v", err)
	}
}

func TestLoadFromBytesRegexValida(t *testing.T) {
	spec, err := LoadFromBytes([]byte(yamlValido))
	if err != nil {
		t.Fatalf("LoadFromBytes retornou erro inesperado: %v", err)
	}
	re, err := spec.Fields[1].Regex()
	if err != nil {
		t.Fatalf("Regex retornou erro inesperado: %v", err)
	}
	if re == nil {
		t.Fatal("Regex deveria devolver a expressão compilada")
	}
	if !re.MatchString("v1.2.3") {
		t.Error("o regex deveria casar com v1.2.3")
	}
	if re.MatchString("1.2.3") {
		t.Error("o regex não deveria casar com 1.2.3")
	}
}

func TestLoadFromBytesRegexInvalida(t *testing.T) {
	_, err := LoadFromBytes([]byte("title: T\nfields:\n  - name: a\n    type: input\n    validate: \"[a-z\"\n"))
	if err == nil {
		t.Fatal("esperado erro para regex inválida, obtido nil")
	}
	if !strings.Contains(err.Error(), "regex de validação inválido") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestLoadFromFileInexistente(t *testing.T) {
	_, err := LoadFromFile("arquivo-que-nao-existe.yaml")
	if err == nil {
		t.Fatal("esperado erro para arquivo inexistente, obtido nil")
	}
	if !strings.Contains(err.Error(), "erro ao ler arquivo") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestNormalizeType(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{"input", TypeInput},
		{"INPUT", TypeInput},
		{" Input ", TypeInput},
		{"text", TypeText},
		{"select", TypeSelect},
		{"confirm", TypeConfirm},
		{"multiselect", TypeMultiSelect},
		{"MultiSelect", TypeMultiSelect},
		{"multi-select", TypeMultiSelect},
	}
	for _, c := range casos {
		obtido, err := NormalizeType(c.entrada)
		if err != nil {
			t.Errorf("NormalizeType(%q) retornou erro: %v", c.entrada, err)
			continue
		}
		if obtido != c.esperado {
			t.Errorf("NormalizeType(%q) = %q, esperado %q", c.entrada, obtido, c.esperado)
		}
	}

	if _, err := NormalizeType("radio"); err == nil {
		t.Error("NormalizeType deveria falhar para tipo desconhecido")
	}
}

func TestFieldSpecLabelText(t *testing.T) {
	f := FieldSpec{Name: "campo"}
	if f.LabelText() != "campo" {
		t.Errorf("LabelText deveria cair para o nome, obtido %q", f.LabelText())
	}
	f.Label = " Rótulo "
	if f.LabelText() != "Rótulo" {
		t.Errorf("LabelText esperado %q, obtido %q", "Rótulo", f.LabelText())
	}
}

// TestExemploValido garante que o formulário de exemplo do repositório
// continua carregando sem erros.
func TestExemploValido(t *testing.T) {
	spec, err := LoadFromFile(filepath.Join("examples", "deploy.yaml"))
	if err != nil {
		t.Fatalf("o exemplo do repositório deveria carregar sem erros: %v", err)
	}
	if len(spec.Fields) == 0 {
		t.Fatal("o exemplo deveria ter campos")
	}
	if spec.Description == "" {
		t.Error("o exemplo deveria ter uma descrição")
	}
}
