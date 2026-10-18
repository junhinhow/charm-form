package charmform

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ToJSON serializa a especificação do formulário para JSON.
//
// O erro é mantido por compatibilidade com a API pública: a especificação só
// contém tipos primitivos, então a serialização na prática não falha.
func ToJSON(spec FormSpec) ([]byte, error) {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar para JSON: %w", err)
	}
	return data, nil
}

// WriteJSON escreve a especificação do formulário em JSON no writer
// informado, terminada por uma quebra de linha.
func WriteJSON(w io.Writer, spec FormSpec) error {
	data, err := ToJSON(spec)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("erro ao escrever JSON: %w", err)
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("erro ao escrever JSON: %w", err)
	}
	return nil
}

// Run imprime a especificação do formulário como JSON no stdout.
//
// Para emitir as respostas coletadas por um formulário use o subpacote
// renderer, ou WriteJSON com um io.Writer próprio.
func Run(spec FormSpec) error {
	return WriteJSON(os.Stdout, spec)
}
