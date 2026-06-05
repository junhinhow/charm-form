package charmform

import (
	"encoding/json"
	"fmt"
)

// ToJSON serializa a especificacao do formulario para JSON.
func ToJSON(spec FormSpec) ([]byte, error) {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar para JSON: %w", err)
	}
	return data, nil
}

// Run imprime a especificacao do formulario como JSON no stdout.
// Placeholder para futura integracao com charmbracelet/huh.
func Run(spec FormSpec) error {
	data, err := ToJSON(spec)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
