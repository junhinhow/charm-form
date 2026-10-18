package charmform

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Answers guarda as respostas de um formulário, indexadas pelo nome do campo.
//
// Os valores usam os seguintes tipos, conforme o tipo do campo:
//
//   - input, text e select: string
//   - confirm: bool
//   - multiselect: []string
type Answers map[string]any

// Validate confere as respostas contra a especificação do formulário.
// Cada campo obrigatório precisa estar preenchido e o regex declarado em
// "validate" precisa casar com o valor informado. Todos os erros são
// reportados de uma vez.
func (a Answers) Validate(spec FormSpec) error {
	var errs []error
	for _, f := range spec.Fields {
		valor, ok := a[f.Name]
		if !ok {
			if f.Required {
				errs = append(errs, fmt.Errorf("campo '%s': obrigatório", f.Name))
			}
			continue
		}
		if err := ValidateField(f, valor); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// ToJSON serializa as respostas em JSON.
func (a Answers) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar respostas para JSON: %w", err)
	}
	return data, nil
}

// WriteJSON escreve as respostas em JSON no writer informado.
func (a Answers) WriteJSON(w io.Writer) error {
	data, err := a.ToJSON()
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("erro ao escrever respostas em JSON: %w", err)
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("erro ao escrever respostas em JSON: %w", err)
	}
	return nil
}

// ValidateField confere a resposta de um único campo contra a especificação
// dele: tipo do valor, obrigatoriedade, opções válidas e regex declarado em
// "validate".
func ValidateField(f FieldSpec, valor any) error {
	switch f.Type {
	case TypeConfirm:
		if _, ok := valor.(bool); !ok {
			return fmt.Errorf("campo '%s': valor deve ser booleano", f.Name)
		}
		// Um booleano, mesmo falso, já é uma resposta válida.
		return nil

	case TypeSelect:
		s, ok := valor.(string)
		if !ok {
			return fmt.Errorf("campo '%s': valor deve ser texto", f.Name)
		}
		if s == "" {
			if f.Required {
				return fmt.Errorf("campo '%s': obrigatório", f.Name)
			}
			return nil
		}
		if !containsOpcao(f.Options, s) {
			return fmt.Errorf("campo '%s': opção '%s' inválida", f.Name, s)
		}
		return validaRegex(f, s)

	case TypeMultiSelect:
		lista, ok := listaDeStrings(valor)
		if !ok {
			return fmt.Errorf("campo '%s': valor deve ser uma lista de textos", f.Name)
		}
		if len(lista) == 0 && f.Required {
			return fmt.Errorf("campo '%s': obrigatório", f.Name)
		}
		for _, s := range lista {
			if !containsOpcao(f.Options, s) {
				return fmt.Errorf("campo '%s': opção '%s' inválida", f.Name, s)
			}
			if err := validaRegex(f, s); err != nil {
				return err
			}
		}
		return nil

	default: // input e text
		s, ok := valor.(string)
		if !ok {
			return fmt.Errorf("campo '%s': valor deve ser texto", f.Name)
		}
		if strings.TrimSpace(s) == "" {
			if f.Required {
				return fmt.Errorf("campo '%s': obrigatório", f.Name)
			}
			return nil
		}
		return validaRegex(f, s)
	}
}

// validaRegex aplica o regex declarado no campo sobre o valor informado.
func validaRegex(f FieldSpec, valor string) error {
	if f.Validate == "" {
		return nil
	}
	re, err := regexp.Compile(f.Validate)
	if err != nil {
		return fmt.Errorf("campo '%s': regex de validação inválido %q: %w", f.Name, f.Validate, err)
	}
	if !re.MatchString(valor) {
		return fmt.Errorf("campo '%s': valor '%s' não corresponde ao padrão %q", f.Name, valor, f.Validate)
	}
	return nil
}

// containsOpcao informa se a opção pertence à lista informada.
func containsOpcao(opcoes []string, alvo string) bool {
	for _, o := range opcoes {
		if o == alvo {
			return true
		}
	}
	return false
}

// listaDeStrings converte o valor informado em uma lista de textos, aceitando
// tanto []string quanto []any (comportamento típico de valores vindos do JSON).
func listaDeStrings(valor any) ([]string, bool) {
	switch v := valor.(type) {
	case nil:
		return nil, true
	case []string:
		return v, true
	case []any:
		lista := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, false
			}
			lista = append(lista, s)
		}
		return lista, true
	default:
		return nil, false
	}
}
