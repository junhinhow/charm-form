package renderer

import "strings"

// textBuffer é um buffer de texto editável com cursor, usado nos campos
// "input" e "text".
type textBuffer struct {
	value  []rune
	cursor int
}

// SetValue substitui o conteúdo do buffer e posiciona o cursor no fim.
func (b *textBuffer) SetValue(s string) {
	b.value = []rune(s)
	b.cursor = len(b.value)
}

// Value devolve o conteúdo do buffer.
func (b *textBuffer) Value() string {
	return string(b.value)
}

// Cursor devolve a posição do cursor no buffer.
func (b *textBuffer) Cursor() int {
	return b.cursor
}

// Insert insere o texto na posição do cursor.
func (b *textBuffer) Insert(s string) {
	r := []rune(s)
	if len(r) == 0 {
		return
	}
	novo := make([]rune, 0, len(b.value)+len(r))
	novo = append(novo, b.value[:b.cursor]...)
	novo = append(novo, r...)
	novo = append(novo, b.value[b.cursor:]...)
	b.value = novo
	b.cursor += len(r)
}

// InsertSemQuebra insere o texto descartando quebras de linha. É usado no
// campo "input", que aceita apenas uma linha.
func (b *textBuffer) InsertSemQuebra(s string) {
	b.Insert(strings.ReplaceAll(s, "\n", " "))
}

// Newline insere uma quebra de linha na posição do cursor.
func (b *textBuffer) Newline() {
	b.Insert("\n")
}

// Backspace apaga o caractere à esquerda do cursor.
func (b *textBuffer) Backspace() {
	if b.cursor == 0 {
		return
	}
	novo := make([]rune, 0, len(b.value)-1)
	novo = append(novo, b.value[:b.cursor-1]...)
	novo = append(novo, b.value[b.cursor:]...)
	b.value = novo
	b.cursor--
}

// Delete apaga o caractere sob o cursor.
func (b *textBuffer) Delete() {
	if b.cursor >= len(b.value) {
		return
	}
	novo := make([]rune, 0, len(b.value)-1)
	novo = append(novo, b.value[:b.cursor]...)
	novo = append(novo, b.value[b.cursor+1:]...)
	b.value = novo
}

// Left move o cursor uma posição para a esquerda.
func (b *textBuffer) Left() {
	if b.cursor > 0 {
		b.cursor--
	}
}

// Right move o cursor uma posição para a direita.
func (b *textBuffer) Right() {
	if b.cursor < len(b.value) {
		b.cursor++
	}
}

// Home move o cursor para o início do buffer.
func (b *textBuffer) Home() {
	b.cursor = 0
}

// End move o cursor para o fim do buffer.
func (b *textBuffer) End() {
	b.cursor = len(b.value)
}

// LinhaDoCursor devolve o índice (base 0) da linha em que o cursor está.
func (b *textBuffer) LinhaDoCursor() int {
	linha := 0
	for _, r := range b.value[:b.cursor] {
		if r == '\n' {
			linha++
		}
	}
	return linha
}

// ColunaDoCursor devolve o índice (base 0) do cursor dentro da linha atual.
func (b *textBuffer) ColunaDoCursor() int {
	coluna := 0
	for _, r := range b.value[:b.cursor] {
		if r == '\n' {
			coluna = 0
			continue
		}
		coluna++
	}
	return coluna
}

// Linhas devolve o conteúdo do buffer dividido em linhas.
func (b *textBuffer) Linhas() []string {
	return strings.Split(string(b.value), "\n")
}
