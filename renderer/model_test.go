package renderer

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	charmform "github.com/junhinhow/charm-form"
)

// campoBase monta um campo de teste com o tipo pedido.
func campoBase(tipo string) charmform.FieldSpec {
	return charmform.FieldSpec{Name: "campo", Type: tipo, Label: "Campo"}
}

// specBase devolve uma especificação simples com um único campo do tipo pedido.
func specBase(tipo string) charmform.FormSpec {
	return charmform.FormSpec{
		Title:       "Título",
		Description: "Descrição",
		Fields:      []charmform.FieldSpec{campoBase(tipo)},
	}
}

// specCom devolve uma especificação com um único campo, já configurado.
func specCom(f charmform.FieldSpec) charmform.FormSpec {
	return charmform.FormSpec{
		Title:       "Título",
		Description: "Descrição",
		Fields:      []charmform.FieldSpec{f},
	}
}

// tecla devolve um KeyPressMsg para um caractere imprimível.
func tecla(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}

// teclaEspecial devolve um KeyPressMsg para uma tecla especial.
func teclaEspecial(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

// teclaMod devolve um KeyPressMsg para uma tecla com modificador.
func teclaMod(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: mod})
}

// digitar envia uma sequência de caracteres para o modelo.
func digitar(m *Model, texto string) {
	for _, r := range texto {
		m.Update(tecla(r))
	}
}

func TestInputDigitaEConfirma(t *testing.T) {
	m := New(charmform.FormSpec{
		Title:  "Título",
		Fields: []charmform.FieldSpec{campoBase(charmform.TypeInput), {Name: "outro", Type: charmform.TypeInput}},
	})

	digitar(m, "abc")
	m.Update(teclaEspecial(tea.KeyEnter))

	respostas := m.Answers()
	valor, ok := respostas["campo"].(string)
	if !ok {
		t.Fatalf("resposta ausente ou com tipo errado: %#v", respostas)
	}
	if valor != "abc" {
		t.Errorf("valor esperado %q, obtido %q", "abc", valor)
	}
	if m.current != 1 {
		t.Errorf("o foco deveria ter avançado para o próximo campo, obtida posição %d", m.current)
	}
	if m.Done() {
		t.Error("o formulário não deveria estar concluído")
	}
}

func TestInputAplicaRegexNaConfirmacao(t *testing.T) {
	f := campoBase(charmform.TypeInput)
	f.Validate = "^[0-9]+$"
	m := New(specCom(f))

	digitar(m, "abc")
	m.Update(teclaEspecial(tea.KeyEnter))

	if _, ok := m.Answers()["campo"]; ok {
		t.Error("o campo inválido não deveria ter sido salvo")
	}
	if m.current != 0 {
		t.Errorf("o foco deveria permanecer no campo, obtida posição %d", m.current)
	}
	if m.fields[0].err == "" {
		t.Error("o erro de validação deveria ser exibido")
	}

	// Um valor válido é aceito depois de limpar o campo.
	for i := 0; i < 3; i++ {
		m.Update(teclaEspecial(tea.KeyBackspace))
	}
	digitar(m, "123")
	m.Update(teclaEspecial(tea.KeyEnter))
	if _, ok := m.Answers()["campo"]; !ok {
		t.Error("o campo válido deveria ter sido salvo")
	}
}

func TestInputObrigatorioVazio(t *testing.T) {
	f := campoBase(charmform.TypeInput)
	f.Required = true
	m := New(specCom(f))

	m.Update(teclaEspecial(tea.KeyEnter))
	if _, ok := m.Answers()["campo"]; ok {
		t.Error("campo obrigatório vazio não deveria ter sido salvo")
	}
	if m.fields[0].err == "" {
		t.Error("o erro de obrigatoriedade deveria ser exibido")
	}
}

func TestSelectNavegaEOpciona(t *testing.T) {
	f := campoBase(charmform.TypeSelect)
	f.Options = []string{"a", "b", "c"}
	m := New(specCom(f))

	m.Update(teclaEspecial(tea.KeyDown))
	m.Update(teclaEspecial(tea.KeyEnter))

	valor, ok := m.Answers()["campo"].(string)
	if !ok {
		t.Fatalf("resposta ausente ou com tipo errado: %#v", m.Answers())
	}
	if valor != "b" {
		t.Errorf("valor esperado %q, obtido %q", "b", valor)
	}
}

func TestSelectNavegaParaCimaNoLimite(t *testing.T) {
	f := campoBase(charmform.TypeSelect)
	f.Options = []string{"a", "b"}
	m := New(specCom(f))

	// No primeiro item, "up" deve dar a volta e ir para o último.
	m.Update(teclaEspecial(tea.KeyUp))
	m.Update(teclaEspecial(tea.KeyEnter))

	valor, _ := m.Answers()["campo"].(string)
	if valor != "b" {
		t.Errorf("valor esperado %q, obtido %q", "b", valor)
	}
}

func TestMultiSelectMarcaOpcoes(t *testing.T) {
	f := campoBase(charmform.TypeMultiSelect)
	f.Options = []string{"a", "b"}
	m := New(specCom(f))

	m.Update(teclaEspecial(tea.KeySpace))
	m.Update(teclaEspecial(tea.KeyDown))
	m.Update(teclaEspecial(tea.KeySpace))
	m.Update(teclaEspecial(tea.KeyEnter))

	valor, ok := m.Answers()["campo"].([]string)
	if !ok {
		t.Fatalf("resposta ausente ou com tipo errado: %#v", m.Answers())
	}
	if len(valor) != 2 || valor[0] != "a" || valor[1] != "b" {
		t.Errorf("valor esperado [a b], obtido %#v", valor)
	}
}

func TestConfirmAlternaEObserve(t *testing.T) {
	m := New(specBase(charmform.TypeConfirm))

	m.Update(tecla('y'))
	m.Update(teclaEspecial(tea.KeyEnter))

	valor, ok := m.Answers()["campo"].(bool)
	if !ok {
		t.Fatalf("resposta ausente ou com tipo errado: %#v", m.Answers())
	}
	if !valor {
		t.Error("o valor esperado era true")
	}
}

func TestTextMultiLinha(t *testing.T) {
	m := New(specBase(charmform.TypeText))

	digitar(m, "linha1")
	m.Update(teclaEspecial(tea.KeyEnter))
	digitar(m, "linha2")
	m.Update(teclaMod('s', tea.ModCtrl))

	valor, ok := m.Answers()["campo"].(string)
	if !ok {
		t.Fatalf("resposta ausente ou com tipo errado: %#v", m.Answers())
	}
	if valor != "linha1\nlinha2" {
		t.Errorf("valor esperado %q, obtido %q", "linha1\nlinha2", valor)
	}
}

func TestCancelamento(t *testing.T) {
	m := New(specBase(charmform.TypeInput))

	_, cmd := m.Update(teclaMod('c', tea.ModCtrl))
	if cmd == nil {
		t.Fatal("ctrl+c deveria encerrar o programa")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Error("ctrl+c deveria devolver o comando de saída")
	}
	if !m.Cancelled() {
		t.Error("o modelo deveria estar marcado como cancelado")
	}
}

func TestNavegacaoPorTab(t *testing.T) {
	m := New(charmform.FormSpec{
		Title: "T",
		Fields: []charmform.FieldSpec{
			{Name: "a", Type: charmform.TypeInput},
			{Name: "b", Type: charmform.TypeInput},
		},
	})

	m.Update(teclaEspecial(tea.KeyTab))
	if m.current != 1 {
		t.Errorf("tab deveria avançar para o campo 2, obtida posição %d", m.current)
	}
	m.Update(teclaMod(tea.KeyTab, tea.ModShift))
	if m.current != 0 {
		t.Errorf("shift+tab deveria voltar para o campo 1, obtida posição %d", m.current)
	}
}

func TestColarTexto(t *testing.T) {
	m := New(specBase(charmform.TypeInput))
	m.Update(tea.PasteMsg{Content: "colado"})
	m.Update(teclaEspecial(tea.KeyEnter))

	valor, _ := m.Answers()["campo"].(string)
	if valor != "colado" {
		t.Errorf("valor esperado %q, obtido %q", "colado", valor)
	}
}

func TestViewUsaTituloEDescricao(t *testing.T) {
	m := New(specBase(charmform.TypeInput))
	conteudo := m.View().Content

	if !strings.Contains(conteudo, "Título") {
		t.Errorf("a view deveria conter o título, obtido:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "Descrição") {
		t.Errorf("a view deveria conter a descrição, obtido:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "Campo") {
		t.Errorf("a view deveria conter o rótulo do campo, obtido:\n%s", conteudo)
	}
}

func TestViewMostraErroDeValidacao(t *testing.T) {
	f := campoBase(charmform.TypeInput)
	f.Validate = "^[0-9]+$"
	m := New(specCom(f))
	digitar(m, "abc")
	m.Update(teclaEspecial(tea.KeyEnter))

	if !strings.Contains(m.View().Content, "não corresponde ao padrão") {
		t.Errorf("a view deveria mostrar o erro de validação, obtido:\n%s", m.View().Content)
	}
}

func TestNewNormalizaTipos(t *testing.T) {
	m := New(charmform.FormSpec{
		Title:  "T",
		Fields: []charmform.FieldSpec{{Name: "a", Type: "Input"}},
	})

	digitar(m, "x")
	m.Update(teclaEspecial(tea.KeyEnter))

	if v, _ := m.Answers()["a"].(string); v != "x" {
		t.Errorf("valor esperado %q, obtido %q", "x", v)
	}
}

func TestTextBuffer(t *testing.T) {
	var b textBuffer
	b.SetValue("abc")
	if b.Value() != "abc" {
		t.Errorf("valor esperado %q, obtido %q", "abc", b.Value())
	}
	if b.Cursor() != 3 {
		t.Errorf("cursor esperado 3, obtido %d", b.Cursor())
	}

	b.Home()
	b.Insert("x")
	if b.Value() != "xabc" {
		t.Errorf("valor esperado %q, obtido %q", "xabc", b.Value())
	}

	b.End()
	b.Backspace()
	if b.Value() != "xab" {
		t.Errorf("valor esperado %q, obtido %q", "xab", b.Value())
	}

	b.Home()
	b.Delete()
	if b.Value() != "ab" {
		t.Errorf("valor esperado %q, obtido %q", "ab", b.Value())
	}

	b.Backspace()
	if b.Value() != "ab" {
		t.Errorf("backspace no início não deveria alterar o valor, obtido %q", b.Value())
	}

	b.Right()
	b.Backspace()
	if b.Value() != "b" {
		t.Errorf("valor esperado %q, obtido %q", "b", b.Value())
	}

	b.InsertSemQuebra("l1\nl2")
	if strings.Contains(b.Value(), "\n") {
		t.Errorf("InsertSemQuebra não deveria inserir quebras de linha, obtido %q", b.Value())
	}

	b.SetValue("a\nb")
	if b.LinhaDoCursor() != 1 {
		t.Errorf("linha do cursor esperada 1, obtida %d", b.LinhaDoCursor())
	}
	if b.ColunaDoCursor() != 1 {
		t.Errorf("coluna do cursor esperada 1, obtida %d", b.ColunaDoCursor())
	}
	if linhas := b.Linhas(); len(linhas) != 2 || linhas[0] != "a" || linhas[1] != "b" {
		t.Errorf("linhas inesperadas: %#v", linhas)
	}
}
