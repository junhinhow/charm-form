// Pacote charmform fornece um gerador declarativo de formularios para TUIs.
// Formularios sao definidos em YAML e podem ser executados como CLIs interativos.
package charmform

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// FormSpec representa a definicao completa de um formulario carregado do YAML.
type FormSpec struct {
	Title       string      `yaml:"title" json:"title"`
	Description string      `yaml:"description" json:"description"`
	Fields      []FieldSpec `yaml:"fields" json:"fields"`
}

// FieldSpec representa um campo individual do formulario.
// Tipos suportados: input, select, confirm, multiselect, text.
type FieldSpec struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type"`
	Label       string   `yaml:"label" json:"label"`
	Placeholder string   `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Options     []string `yaml:"options,omitempty" json:"options,omitempty"`
	Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
	Validate    string   `yaml:"validate,omitempty" json:"validate,omitempty"`
}

// LoadFromFile carrega uma especificacao de formulario a partir de um arquivo YAML.
func LoadFromFile(path string) (FormSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FormSpec{}, fmt.Errorf("erro ao ler arquivo %s: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes carrega uma especificacao de formulario a partir de bytes YAML.
func LoadFromBytes(data []byte) (FormSpec, error) {
	var spec FormSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return FormSpec{}, fmt.Errorf("erro ao parsear YAML: %w", err)
	}
	if spec.Title == "" {
		return FormSpec{}, fmt.Errorf("campo 'title' e obrigatorio na especificacao do formulario")
	}
	if len(spec.Fields) == 0 {
		return FormSpec{}, fmt.Errorf("formulario deve ter pelo menos um campo")
	}
	for i, f := range spec.Fields {
		if f.Name == "" {
			return FormSpec{}, fmt.Errorf("campo %d: 'name' e obrigatorio", i)
		}
		if f.Type == "" {
			return FormSpec{}, fmt.Errorf("campo '%s': 'type' e obrigatorio", f.Name)
		}
		switch f.Type {
		case "input", "select", "confirm", "multiselect", "text":
			// tipos validos
		default:
			return FormSpec{}, fmt.Errorf("campo '%s': tipo '%s' invalido (use input, select, confirm, multiselect ou text)", f.Name, f.Type)
		}
		if (f.Type == "select" || f.Type == "multiselect") && len(f.Options) == 0 {
			return FormSpec{}, fmt.Errorf("campo '%s': tipo '%s' requer pelo menos uma opcao", f.Name, f.Type)
		}
	}
	return spec, nil
}
