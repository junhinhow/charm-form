# charm-form

> Gerador declarativo de formulários para TUIs -- defina formulários em YAML, obtenha apps CLI interativos.

[Read in English](README.md)

## Problema

Criar formulários interativos para CLI exige escrever código Go com bibliotecas TUI específicas. Cada novo formulário significa mais boilerplate, mais ciclos de compilação e conhecimento em Go só para fazer algumas perguntas.

## Solução

**charm-form** permite definir formulários em YAML puro e obter uma interface de terminal totalmente interativa. Nenhum conhecimento em Go é necessário para criar ou modificar formulários -- basta editar um arquivo YAML.

O módulo traz três partes:

| Parte | Caminho de importação | O que faz |
|---|---|---|
| Parser e validador da especificação | `github.com/junhinhow/charm-form` | Carrega o YAML, valida e aplica a regex de `validate` |
| Renderizador interativo | `github.com/junhinhow/charm-form/renderer` | Modelo Bubble Tea v2 que desenha os campos e coleta as respostas |
| CLI | `github.com/junhinhow/charm-form/cmd/charm-form` | Executa um formulário e grava as respostas em JSON |

## Instalação

Instale a CLI:

```bash
go install github.com/junhinhow/charm-form/cmd/charm-form@latest
```

Ou adicione como biblioteca:

```bash
go get github.com/junhinhow/charm-form
```

> Requer Go 1.25 ou superior -- o piso imposto pelo `charm.land/lipgloss/v2`.

## Início rápido

```bash
charm-form -file examples/deploy.yaml
```

A CLI abre o formulário no terminal, coleta as respostas e as imprime em JSON:

```json
{
  "environment": "production",
  "version": "v1.4.0",
  "dry_run": false,
  "services": ["api", "web"]
}
```

## Referência da Especificação YAML

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

| Tipo | Descrição | Tipo da resposta |
|------|-----------|------------------|
| `input` | Entrada de texto em linha unica | `string` |
| `text` | Entrada de texto em multiplas linhas | `string` |
| `select` | Escolha unica entre opcoes | `string` |
| `multiselect` | Multiplas escolhas entre opcoes | `[]string` |
| `confirm` | Booleano sim/nao | `bool` |

### Regras de validação

- `required: true` -- o campo não pode ficar vazio.
- `validate: "<regex>"` -- a resposta precisa casar com a expressão regular. O padrão é compilado durante o carregamento da especificação, então uma regex inválida é reportada na hora em vez de falhar em silêncio.
- `name` precisa ser único no formulário; nomes duplicados são rejeitados durante o carregamento da especificação.
- `type` não diferencia maiúsculas de minúsculas (`Input`, `INPUT` e `input` são o mesmo tipo) e é normalizado para minúsculas.

## CLI

```
Uso: charm-form [opcoes] [formulario.yaml]

  -file string
        caminho do arquivo YAML com a definicao do formulario
  -out string
        arquivo onde o resultado e gravado (padrao: stdout)
  -json bool
        emite as respostas em JSON; use -json=false para um resumo legivel
```

## Teclas

| Tecla | Ação |
|-------|------|
| `enter` | Confirma o campo atual |
| `tab` / `shift+tab` | Próximo / campo anterior |
| `seta para cima` / `seta para baixo` | Move a seleção (`select` e `multiselect`) |
| `espaço` | Marca ou desmarca uma opção (`multiselect`) |
| `y` / `n` | Alterna um booleano (`confirm`) |
| `ctrl+s` | Confirma um campo `text` de múltiplas linhas (`enter` insere uma quebra de linha) |
| `ctrl+c` / `esc` | Cancela o formulário |

## Uso como biblioteca

```go
package main

import (
    "log"
    "os"

    charmform "github.com/junhinhow/charm-form"
    "github.com/junhinhow/charm-form/renderer"
)

func main() {
    spec, err := charmform.LoadFromFile("deploy.yaml")
    if err != nil {
        log.Fatal(err)
    }

    answers, err := renderer.Run(spec)
    if err != nil {
        log.Fatal(err)
    }

    if err := answers.WriteJSON(os.Stdout); err != nil {
        log.Fatal(err)
    }
}
```

### Trabalhando apenas com a especificação

`charmform.Run(spec)` imprime a própria especificação em JSON -- útil para inspecionar ou converter uma definição de formulário:

```go
spec, err := charmform.LoadFromFile("deploy.yaml")
if err != nil {
    log.Fatal(err)
}
if err := charmform.Run(spec); err != nil {
    log.Fatal(err)
}
```

`charmform.WriteJSON(w, spec)` faz o mesmo em qualquer `io.Writer`.

## Exemplos

Veja o diretório [examples/](examples/) para exemplos de definições de formulários YAML.

## Roadmap

- [x] Renderização TUI interativa com Bubble Tea v2
- [x] Ferramenta CLI (`charm-form -file deploy.yaml`)
- [x] Saída em JSON
- [x] Validação por regex (`validate`)
- [ ] Formatos de saída: variáveis de ambiente, TOML
- [ ] Campos condicionais (mostrar/ocultar baseado em outras respostas)
- [ ] Valores padrão e interpolação de variáveis de ambiente
- [ ] Composição de formulários (incluir/estender outros YAML)
- [ ] Validação além de regex (min/max, email, URL, etc.)

## Licença

MIT
