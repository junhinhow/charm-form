package charmform

import (
	"strings"
	"testing"
)

func TestValidateFieldRegex(t *testing.T) {
	f := FieldSpec{Name: "version", Type: TypeInput, Validate: "^v\\d+\\.\\d+\\.\\d+$"}

	if err := ValidateField(f, "v1.2.3"); err != nil {
		t.Errorf("valor válido rejeitado: %v", err)
	}
	err := ValidateField(f, "1.2.3")
	if err == nil {
		t.Fatal("valor inválido aceito")
	}
	if !strings.Contains(err.Error(), "não corresponde ao padrão") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

func TestValidateFieldRegexEmMultiselect(t *testing.T) {
	f := FieldSpec{
		Name:     "tags",
		Type:     TypeMultiSelect,
		Options:  []string{"v1", "v2"},
		Validate: "^v\\d+$",
	}

	if err := ValidateField(f, []string{"v1", "v2"}); err != nil {
		t.Errorf("lista válida rejeitada: %v", err)
	}
	if err := ValidateField(f, []string{"v1", "x"}); err == nil {
		t.Error("lista com item inválido foi aceita")
	}
}

func TestValidateFieldObrigatorio(t *testing.T) {
	f := FieldSpec{Name: "nome", Type: TypeInput, Required: true}

	if err := ValidateField(f, ""); err == nil {
		t.Error("campo obrigatório vazio foi aceito")
	}
	if err := ValidateField(f, "   "); err == nil {
		t.Error("campo obrigatório com espaços foi aceito")
	}
	if err := ValidateField(f, "valor"); err != nil {
		t.Errorf("campo obrigatório preenchido foi rejeitado: %v", err)
	}
}

func TestValidateFieldObrigatorioMultiselect(t *testing.T) {
	f := FieldSpec{Name: "servicos", Type: TypeMultiSelect, Options: []string{"api"}, Required: true}

	if err := ValidateField(f, []string{}); err == nil {
		t.Error("multiselect obrigatório vazio foi aceito")
	}
	if err := ValidateField(f, []string{"api"}); err != nil {
		t.Errorf("multiselect obrigatório preenchido foi rejeitado: %v", err)
	}
}

func TestValidateFieldConfirm(t *testing.T) {
	f := FieldSpec{Name: "dry_run", Type: TypeConfirm}

	if err := ValidateField(f, false); err != nil {
		t.Errorf("confirm falso foi rejeitado: %v", err)
	}
	if err := ValidateField(f, "não"); err == nil {
		t.Error("confirm com valor textual deveria ser rejeitado")
	}
}

func TestValidateFieldOpcaoInvalida(t *testing.T) {
	f := FieldSpec{Name: "ambiente", Type: TypeSelect, Options: []string{"dev", "prod"}}

	if err := ValidateField(f, "prod"); err != nil {
		t.Errorf("opção válida rejeitada: %v", err)
	}
	if err := ValidateField(f, "qa"); err == nil {
		t.Error("opção inválida foi aceita")
	}
}

func TestAnswersValidate(t *testing.T) {
	spec := FormSpec{
		Title: "T",
		Fields: []FieldSpec{
			{Name: "nome", Type: TypeInput, Required: true},
			{Name: "email", Type: TypeInput, Validate: "^[^@]+@[^@]+$"},
			{Name: "ativo", Type: TypeConfirm},
		},
	}

	answers := Answers{
		"nome":  "Ana",
		"email": "ana@example.com",
		"ativo": true,
	}
	if err := answers.Validate(spec); err != nil {
		t.Errorf("respostas válidas rejeitadas: %v", err)
	}

	answers["email"] = "invalido"
	err := answers.Validate(spec)
	if err == nil {
		t.Fatal("email inválido foi aceito")
	}
	if !strings.Contains(err.Error(), "email") {
		t.Errorf("o erro deveria citar o campo, obtido: %v", err)
	}
}

func TestAnswersValidateRelatorioTodosOsErros(t *testing.T) {
	spec := FormSpec{
		Title: "T",
		Fields: []FieldSpec{
			{Name: "nome", Type: TypeInput, Required: true},
			{Name: "email", Type: TypeInput, Validate: "^[^@]+@[^@]+$"},
		},
	}

	err := Answers{"email": "invalido"}.Validate(spec)
	if err == nil {
		t.Fatal("email inválido foi aceito")
	}
	texto := err.Error()
	if !strings.Contains(texto, "nome") || !strings.Contains(texto, "email") {
		t.Errorf("o erro deveria listar os dois campos, obtido: %v", err)
	}
}

func TestAnswersToJSON(t *testing.T) {
	answers := Answers{"nome": "Ana", "ativo": true, "tags": []string{"a", "b"}}
	data, err := answers.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON retornou erro inesperado: %v", err)
	}
	texto := string(data)
	for _, esperado := range []string{`"nome": "Ana"`, `"ativo": true`, `"tags"`} {
		if !strings.Contains(texto, esperado) {
			t.Errorf("JSON deveria conter %s, obtido:\n%s", esperado, texto)
		}
	}
}

func TestAnswersWriteJSON(t *testing.T) {
	var b strings.Builder
	if err := (Answers{"nome": "Ana"}).WriteJSON(&b); err != nil {
		t.Fatalf("WriteJSON retornou erro inesperado: %v", err)
	}
	if !strings.HasSuffix(b.String(), "\n") {
		t.Error("WriteJSON deveria terminar com quebra de linha")
	}
	if !strings.Contains(b.String(), `"nome": "Ana"`) {
		t.Errorf("conteúdo inesperado: %s", b.String())
	}
}
