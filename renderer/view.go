package renderer

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	charmform "github.com/junhinhow/charm-form"
)

// Estilos usados na renderização do formulário.
var (
	estiloTitulo      = lipgloss.NewStyle().Bold(true)
	estiloDescricao   = lipgloss.NewStyle().Faint(true)
	estiloRotulo      = lipgloss.NewStyle().Bold(true)
	estiloObrigatorio = lipgloss.NewStyle().Faint(true)
	estiloErro        = lipgloss.NewStyle().Bold(true)
	estiloAjuda       = lipgloss.NewStyle().Faint(true)
	estiloSelecao     = lipgloss.NewStyle().Bold(true)
	estiloCursor      = lipgloss.NewStyle().Reverse(true)
	estiloPlaceholder = lipgloss.NewStyle().Faint(true)
	estiloContador    = lipgloss.NewStyle().Faint(true)
)

// render monta a visão completa do formulário.
func (m *Model) render() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render(m.spec.Title))
	b.WriteString("\n")
	if strings.TrimSpace(m.spec.Description) != "" {
		b.WriteString(estiloDescricao.Render(m.spec.Description))
		b.WriteString("\n")
	}
	b.WriteString(estiloContador.Render(fmt.Sprintf("campo %d de %d", m.current+1, len(m.fields))))
	b.WriteString("\n\n")

	for i, c := range m.fields {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(m.renderCampo(c, i == m.current))
	}

	b.WriteString("\n\n")
	b.WriteString(estiloAjuda.Render(m.ajuda()))
	b.WriteString("\n")

	return m.aplicaLargura(b.String())
}

// renderCampo renderiza um campo do formulário.
func (m *Model) renderCampo(c *campo, ativo bool) string {
	var b strings.Builder

	prefixo := "  "
	if ativo {
		prefixo = "  " + estiloSelecao.Render("▸ ")
	}

	rotulo := c.spec.LabelText()
	if c.spec.Required {
		rotulo += estiloObrigatorio.Render(" *")
	}
	b.WriteString(prefixo)
	b.WriteString(estiloRotulo.Render(rotulo))
	b.WriteString("\n")

	switch c.spec.Type {
	case charmform.TypeInput, charmform.TypeText:
		b.WriteString(m.renderEntrada(c, ativo))
	case charmform.TypeSelect:
		b.WriteString(m.renderSelect(c, ativo))
	case charmform.TypeConfirm:
		b.WriteString(m.renderConfirm(c, ativo))
	case charmform.TypeMultiSelect:
		b.WriteString(m.renderMulti(c, ativo))
	}

	if c.err != "" {
		b.WriteString("\n")
		b.WriteString(estiloErro.Render("  ! " + c.err))
	}
	return b.String()
}

// renderEntrada renderiza os campos "input" e "text".
func (m *Model) renderEntrada(c *campo, ativo bool) string {
	linhas := c.texto.Linhas()
	linhaCursor := c.texto.LinhaDoCursor()
	vazio := c.texto.Value() == ""

	var out []string
	for i, linha := range linhas {
		texto := linha
		switch {
		case ativo && i == linhaCursor:
			if vazio && c.spec.Placeholder != "" {
				texto = estiloCursor.Render(" ") + estiloPlaceholder.Render(c.spec.Placeholder)
				break
			}
			runes := []rune(linha)
			coluna := c.texto.ColunaDoCursor()
			if coluna > len(runes) {
				coluna = len(runes)
			}
			texto = string(runes[:coluna]) + estiloCursor.Render(" ") + string(runes[coluna:])
		case vazio && c.spec.Placeholder != "" && i == 0:
			texto = estiloPlaceholder.Render(c.spec.Placeholder)
		}
		out = append(out, "    "+texto)
	}
	return strings.Join(out, "\n")
}

// renderSelect renderiza o campo "select".
func (m *Model) renderSelect(c *campo, ativo bool) string {
	var b strings.Builder
	for i, o := range c.spec.Options {
		marca := "  "
		estilo := lipgloss.NewStyle()
		if ativo && i == c.cursor {
			marca = estiloSelecao.Render("❯ ")
			estilo = estiloSelecao
		}
		b.WriteString("    " + marca + estilo.Render(o) + "\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// renderMulti renderiza o campo "multiselect".
func (m *Model) renderMulti(c *campo, ativo bool) string {
	var b strings.Builder
	for i, o := range c.spec.Options {
		marca := "[ ]"
		if c.marcado[i] {
			marca = "[x]"
		}
		prefixo := "  "
		estilo := lipgloss.NewStyle()
		if ativo && i == c.cursor {
			prefixo = estiloSelecao.Render("❯ ")
			estilo = estiloSelecao
		}
		b.WriteString("    " + prefixo + marca + " " + estilo.Render(o) + "\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// renderConfirm renderiza o campo "confirm".
func (m *Model) renderConfirm(c *campo, ativo bool) string {
	sim, nao := "( )", "( )"
	if c.confirm {
		sim = "(●)"
	} else {
		nao = "(●)"
	}
	linha := "    " + sim + " Sim    " + nao + " Não"
	if ativo {
		linha = estiloSelecao.Render(linha)
	}
	return linha
}

// ajuda devolve a linha de ajuda das teclas para o campo corrente.
func (m *Model) ajuda() string {
	if len(m.fields) == 0 {
		return "esc cancelar"
	}
	switch m.fields[m.current].spec.Type {
	case charmform.TypeText:
		return "ctrl+s confirmar  •  enter quebra de linha  •  tab próximo  •  shift+tab anterior  •  esc cancelar"
	case charmform.TypeSelect:
		return "↑/↓ escolher  •  enter confirmar  •  tab próximo  •  shift+tab anterior  •  esc cancelar"
	case charmform.TypeConfirm:
		return "y/n alternar  •  enter confirmar  •  tab próximo  •  shift+tab anterior  •  esc cancelar"
	case charmform.TypeMultiSelect:
		return "espaço marcar  •  ↑/↓ mover  •  enter confirmar  •  tab próximo  •  shift+tab anterior  •  esc cancelar"
	default:
		return "enter confirmar  •  tab próximo  •  shift+tab anterior  •  esc cancelar"
	}
}

// aplicaLargura limita a largura da visão à largura do terminal, quando ela
// é conhecida.
func (m *Model) aplicaLargura(s string) string {
	if m.width <= 0 {
		return s
	}
	return lipgloss.NewStyle().MaxWidth(m.width).Render(s)
}
