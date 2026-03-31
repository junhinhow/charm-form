# charm-form

> Gerador declarativo de formularios para TUIs -- defina formularios em YAML, obtenha apps CLI interativos.

[Read in English](README.md)

## Problema

Criar formularios interativos para CLI exige escrever codigo Go com bibliotecas TUI especificas. Cada novo formulario significa mais boilerplate, mais ciclos de compilacao e conhecimento em Go so para fazer algumas perguntas.

## Solucao

**charm-form** permite definir formularios em YAML puro e obter uma interface de terminal totalmente interativa. Nenhum conhecimento em Go e necessario para criar ou modificar formularios -- basta editar um arquivo YAML.

## Instalacao

```bash
go install github.com/junhinhow/charm-form@latest
```

Ou adicione como biblioteca:

```bash
go get github.com/junhinhow/charm-form
```

## Referencia da Especificacao YAML

```yaml
title: Meu Formulario       # Obrigatorio - titulo do formulario
description: Alguma info     # Opcional - exibido abaixo do titulo
fields:                      # Obrigatorio - pelo menos um campo
  - name: id_campo           # Obrigatorio - identificador unico
    type: input              # Obrigatorio - input|select|confirm|multiselect|text
    label: Label Exibido     # Label do campo mostrado ao usuario
    placeholder: dica        # Texto placeholder (apenas input/text)
    options:                 # Obrigatorio para select/multiselect
      - opcao1
      - opcao2
    required: true           # Se o campo e obrigatorio
    validate: "^regex$"      # Padrao regex de validacao (opcional)
```

### Tipos de Campo

| Tipo | Descricao |
|------|-----------|
| `input` | Entrada de texto em linha unica |
| `text` | Entrada de texto em multiplas linhas |
| `select` | Escolha unica entre opcoes |
| `multiselect` | Multiplas escolhas entre opcoes |
| `confirm` | Booleano sim/nao |

## Uso

```go
package main

import (
    "log"
    charmform "github.com/junhinhow/charm-form"
)

func main() {
    spec, err := charmform.LoadFromFile("deploy.yaml")
    if err != nil {
        log.Fatal(err)
    }
    if err := charmform.Run(spec); err != nil {
        log.Fatal(err)
    }
}
```

## Exemplos

Veja o diretorio [examples/](examples/) para exemplos de definicoes de formularios YAML.

## Roadmap

- [ ] Renderizacao TUI interativa com [charmbracelet/huh](https://github.com/charmbracelet/huh)
- [ ] Ferramenta CLI (`charm-form run deploy.yaml`)
- [ ] Formatos de saida: JSON, variaveis de ambiente, TOML
- [ ] Campos condicionais (mostrar/ocultar baseado em outras respostas)
- [ ] Valores padrao e interpolacao de variaveis de ambiente
- [ ] Composicao de formularios (incluir/estender outros YAML)
- [ ] Validacao alem de regex (min/max, email, URL, etc.)

## Licenca

MIT
