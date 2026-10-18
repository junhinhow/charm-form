// Package renderer implementa a interface de terminal interativa (TUI) do
// charm-form, baseada em Bubble Tea v2 e Lip Gloss v2.
package renderer

import (
	"errors"
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"

	charmform "github.com/junhinhow/charm-form"
)

// ErrCancelado é devolvido por Run quando o usuário cancela o formulário.
var ErrCancelado = errors.New("formulário cancelado pelo usuário")

// Model é o modelo Bubble Tea do formulário interativo.
type Model struct {
	spec    charmform.FormSpec
	fields  []*campo
	current int
	width   int
	height  int
	answers charmform.Answers
	cancel  bool
	done    bool
}

// campo guarda o estado editável de um campo do formulário.
type campo struct {
	spec    charmform.FieldSpec
	texto   textBuffer
	cursor  int
	marcado map[int]bool
	confirm bool
	err     string
}

// New cria um modelo para a especificação informada.
func New(spec charmform.FormSpec) *Model {
	m := &Model{
		spec:    spec,
		answers: charmform.Answers{},
	}
	for _, f := range spec.Fields {
		m.fields = append(m.fields, &campo{spec: f, marcado: map[int]bool{}})
	}
	return m
}

// Answers devolve uma cópia das respostas coletadas até o momento.
func (m *Model) Answers() charmform.Answers {
	out := make(charmform.Answers, len(m.answers))
	for k, v := range m.answers {
		out[k] = v
	}
	return out
}

// Cancelled informa se o usuário cancelou o formulário.
func (m *Model) Cancelled() bool {
	return m.cancel
}

// Done informa se o formulário foi concluído.
func (m *Model) Done() bool {
	return m.done
}

// Init implementa tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implementa tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.PasteMsg:
		m.colaTexto(msg.Content)
		return m, nil
	case tea.KeyPressMsg:
		return m.trataTecla(msg)
	}
	return m, nil
}

// View implementa tea.Model.
func (m *Model) View() tea.View {
	return tea.NewView(m.render())
}

// trataTecla despacha a tecla pressionada para o campo corrente.
func (m *Model) trataTecla(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c", "esc":
		m.cancel = true
		return m, tea.Quit
	}

	if len(m.fields) == 0 {
		return m, nil
	}
	c := m.fields[m.current]

	switch c.spec.Type {
	case charmform.TypeInput, charmform.TypeText:
		return m.trataTeclaTexto(c, key)
	case charmform.TypeSelect:
		return m.trataTeclaSelect(c, key)
	case charmform.TypeConfirm:
		return m.trataTeclaConfirm(c, key)
	case charmform.TypeMultiSelect:
		return m.trataTeclaMulti(c, key)
	}
	return m, nil
}

// trataTeclaTexto trata as teclas dos campos "input" e "text".
func (m *Model) trataTeclaTexto(c *campo, key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "enter":
		// Em campos de múltiplas linhas, "enter" insere uma quebra de linha
		// e "ctrl+s" confirma o valor.
		if c.spec.Type == charmform.TypeText {
			c.texto.Newline()
			c.err = ""
			return m, nil
		}
		return m.confirma(c)
	case "ctrl+s":
		return m.confirma(c)
	case "tab":
		return m.proximo()
	case "shift+tab":
		return m.anterior()
	case "backspace":
		c.texto.Backspace()
	case "delete":
		c.texto.Delete()
	case "left":
		c.texto.Left()
	case "right":
		c.texto.Right()
	case "home":
		c.texto.Home()
	case "end":
		c.texto.End()
	default:
		if key.Text != "" && !key.Mod.Contains(tea.ModCtrl) && !key.Mod.Contains(tea.ModAlt) {
			if c.spec.Type == charmform.TypeInput {
				c.texto.InsertSemQuebra(key.Text)
			} else {
				c.texto.Insert(key.Text)
			}
		}
	}
	c.err = ""
	return m, nil
}

// trataTeclaSelect trata as teclas do campo "select".
func (m *Model) trataTeclaSelect(c *campo, key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	total := len(c.spec.Options)
	if total == 0 {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		c.cursor = (c.cursor - 1 + total) % total
	case "down", "j":
		c.cursor = (c.cursor + 1) % total
	case "enter":
		return m.confirma(c)
	case "tab":
		return m.proximo()
	case "shift+tab":
		return m.anterior()
	}
	c.err = ""
	return m, nil
}

// trataTeclaConfirm trata as teclas do campo "confirm".
func (m *Model) trataTeclaConfirm(c *campo, key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "left", "right", "y", "Y":
		c.confirm = !c.confirm
	case "n", "N":
		c.confirm = false
	case "enter":
		return m.confirma(c)
	case "tab":
		return m.proximo()
	case "shift+tab":
		return m.anterior()
	}
	c.err = ""
	return m, nil
}

// trataTeclaMulti trata as teclas do campo "multiselect".
func (m *Model) trataTeclaMulti(c *campo, key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	total := len(c.spec.Options)
	if total == 0 {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		c.cursor = (c.cursor - 1 + total) % total
	case "down", "j":
		c.cursor = (c.cursor + 1) % total
	case "space":
		c.marcado[c.cursor] = !c.marcado[c.cursor]
	case "enter":
		return m.confirma(c)
	case "tab":
		return m.proximo()
	case "shift+tab":
		return m.anterior()
	}
	c.err = ""
	return m, nil
}

// colaTexto insere texto colado (bracketed paste) no campo corrente.
func (m *Model) colaTexto(conteudo string) {
	if len(m.fields) == 0 || conteudo == "" {
		return
	}
	c := m.fields[m.current]
	switch c.spec.Type {
	case charmform.TypeInput:
		c.texto.InsertSemQuebra(conteudo)
	case charmform.TypeText:
		c.texto.Insert(conteudo)
	}
	c.err = ""
}

// confirma valida o valor do campo corrente. Se o valor for válido, ele é
// guardado e o foco avança para o próximo campo; caso contrário, o erro é
// exibido e o foco permanece no campo.
func (m *Model) confirma(c *campo) (tea.Model, tea.Cmd) {
	valor := c.valor()
	if err := charmform.ValidateField(c.spec, valor); err != nil {
		c.err = err.Error()
		return m, nil
	}
	c.err = ""
	m.answers[c.spec.Name] = valor

	if m.current >= len(m.fields)-1 {
		m.done = true
		return m, tea.Quit
	}
	m.current++
	return m, nil
}

// proximo move o foco para o próximo campo.
func (m *Model) proximo() (tea.Model, tea.Cmd) {
	if m.current < len(m.fields)-1 {
		m.current++
	}
	return m, nil
}

// anterior move o foco para o campo anterior.
func (m *Model) anterior() (tea.Model, tea.Cmd) {
	if m.current > 0 {
		m.current--
	}
	return m, nil
}

// valor devolve o valor atual do campo, conforme o seu tipo.
func (c *campo) valor() any {
	switch c.spec.Type {
	case charmform.TypeConfirm:
		return c.confirm
	case charmform.TypeSelect:
		if c.cursor < 0 || c.cursor >= len(c.spec.Options) {
			return ""
		}
		return c.spec.Options[c.cursor]
	case charmform.TypeMultiSelect:
		lista := make([]string, 0, len(c.spec.Options))
		for i, o := range c.spec.Options {
			if c.marcado[i] {
				lista = append(lista, o)
			}
		}
		return lista
	default:
		return c.texto.Value()
	}
}

// Options configura a execução do formulário.
type Options struct {
	// Input é o leitor de onde as teclas chegam. Quando nil, o programa usa
	// o terminal.
	Input io.Reader
	// Output é o writer onde a interface é desenhada. Quando nil, o programa
	// usa o terminal.
	Output io.Writer
	// Width e Height informam o tamanho da janela. Útil em testes, onde não
	// existe terminal real.
	Width  int
	Height int
}

// Run executa o formulário interativo e devolve as respostas coletadas.
// Devolve ErrCancelado quando o usuário cancela o formulário.
func Run(spec charmform.FormSpec) (charmform.Answers, error) {
	return RunWithOptions(spec, Options{})
}

// RunWithOptions executa o formulário interativo com as opções informadas e
// devolve as respostas coletadas. Devolve ErrCancelado quando o usuário
// cancela o formulário.
func RunWithOptions(spec charmform.FormSpec, opts Options) (charmform.Answers, error) {
	m := New(spec)

	var programOptions []tea.ProgramOption
	if opts.Input != nil {
		programOptions = append(programOptions, tea.WithInput(opts.Input))
	}
	if opts.Output != nil {
		programOptions = append(programOptions, tea.WithOutput(opts.Output))
	}
	if opts.Width > 0 && opts.Height > 0 {
		programOptions = append(programOptions, tea.WithWindowSize(opts.Width, opts.Height))
	}

	programa := tea.NewProgram(m, programOptions...)
	final, err := programa.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar o formulário: %w", err)
	}
	fm, ok := final.(*Model)
	if !ok {
		return nil, fmt.Errorf("modelo de formulário inesperado (%T)", final)
	}
	if fm.cancel {
		return nil, ErrCancelado
	}
	return fm.Answers(), nil
}
