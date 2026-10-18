// Package charmform fornece um gerador declarativo de formulários para TUIs.
// Formulários são definidos em YAML e podem ser executados como CLIs
// interativos, usando o renderizador do subpacote renderer.
package charmform

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Tipos de campo suportados.
const (
	// TypeInput é um campo de texto de linha única.
	TypeInput = "input"
	// TypeText é um campo de texto de múltiplas linhas.
	TypeText = "text"
	// TypeSelect é um campo de escolha única.
	TypeSelect = "select"
	// TypeMultiSelect é um campo de escolha múltipla.
	TypeMultiSelect = "multiselect"
	// TypeConfirm é um campo booleano (sim/não).
	TypeConfirm = "confirm"
)

// FormSpec representa a definição completa de um formulário carregado do YAML.
type FormSpec struct {
	Title       string      `yaml:"title" json:"title"`
	Description string      `yaml:"description" json:"description"`
	Fields      []FieldSpec `yaml:"fields" json:"fields"`
}

// FieldSpec representa um campo individual do formulário.
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

// LoadFromFile carrega uma especificação de formulário a partir de um arquivo YAML.
func LoadFromFile(path string) (FormSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FormSpec{}, fmt.Errorf("erro ao ler arquivo %s: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes carrega uma especificação de formulário a partir de bytes YAML.
// A especificação é validada e os tipos de campo são normalizados para a
// forma canérica (minúscula).
func LoadFromBytes(data []byte) (FormSpec, error) {
	var spec FormSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return FormSpec{}, fmt.Errorf("erro ao parsear YAML: %w", err)
	}
	spec = spec.Normalize()
	if err := spec.Validate(); err != nil {
		return FormSpec{}, err
	}
	return spec, nil
}

// Normalize devolve uma cópia da especificação com os tipos de campo
// normalizados para a forma canônica. Tipos desconhecidos são preservados
// para que a validação possa reportar o erro.
func (s FormSpec) Normalize() FormSpec {
	out := s
	out.Fields = make([]FieldSpec, len(s.Fields))
	copy(out.Fields, s.Fields)
	for i := range out.Fields {
		if t, err := NormalizeType(out.Fields[i].Type); err == nil {
			out.Fields[i].Type = t
		}
	}
	return out
}

// Validate confere a especificação do formulário: título, campos, tipos
// válidos, unicidade dos nomes, opções de seleção e regex de validação.
// Os índices de campo informados nos erros são 1-based, como o usuário conta.
func (s FormSpec) Validate() error {
	if strings.TrimSpace(s.Title) == "" {
		return errors.New("campo 'title' é obrigatório na especificação do formulário")
	}
	if len(s.Fields) == 0 {
		return errors.New("formulário deve ter pelo menos um campo")
	}

	vistos := make(map[string]int, len(s.Fields))
	for i := range s.Fields {
		f := s.Fields[i]
		pos := i + 1

		if strings.TrimSpace(f.Name) == "" {
			return fmt.Errorf("campo %d: 'name' é obrigatório", pos)
		}
		if anterior, dup := vistos[f.Name]; dup {
			return fmt.Errorf("campo %d: nome duplicado '%s' (já usado no campo %d)", pos, f.Name, anterior)
		}
		vistos[f.Name] = pos

		if strings.TrimSpace(f.Type) == "" {
			return fmt.Errorf("campo '%s': 'type' é obrigatório", f.Name)
		}
		tipo, err := NormalizeType(f.Type)
		if err != nil {
			return fmt.Errorf("campo '%s': %w", f.Name, err)
		}
		if (tipo == TypeSelect || tipo == TypeMultiSelect) && len(f.Options) == 0 {
			return fmt.Errorf("campo '%s': tipo '%s' requer pelo menos uma opção", f.Name, tipo)
		}
		if f.Validate != "" {
			if _, err := regexp.Compile(f.Validate); err != nil {
				return fmt.Errorf("campo '%s': regex de validação inválido %q: %w", f.Name, f.Validate, err)
			}
		}
	}
	return nil
}

// NormalizeType normaliza o tipo de um campo. A comparação não diferencia
// maiúsculas de minúsculas e aceita variações comuns de escrita
// ("Input", "MULTISELECT", "multi-select").
func NormalizeType(tipo string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(tipo)) {
	case TypeInput:
		return TypeInput, nil
	case TypeText:
		return TypeText, nil
	case TypeSelect:
		return TypeSelect, nil
	case TypeConfirm:
		return TypeConfirm, nil
	case TypeMultiSelect, "multi-select", "multi_select":
		return TypeMultiSelect, nil
	default:
		return "", fmt.Errorf("tipo '%s' inválido (use input, select, confirm, multiselect ou text)", tipo)
	}
}

// Label devolve o texto exibido ao usuário para o campo, caindo para o
// identificador quando o YAML não define um rótulo.
func (f FieldSpec) LabelText() string {
	if rotulo := strings.TrimSpace(f.Label); rotulo != "" {
		return rotulo
	}
	return f.Name
}

// Regex compila a expressão regular declarada em "validate". Quando o campo
// não declara regex, devolve nil e nenhum erro.
func (f FieldSpec) Regex() (*regexp.Regexp, error) {
	if f.Validate == "" {
		return nil, nil
	}
	re, err := regexp.Compile(f.Validate)
	if err != nil {
		return nil, fmt.Errorf("regex de validação inválido %q: %w", f.Validate, err)
	}
	return re, nil
}
