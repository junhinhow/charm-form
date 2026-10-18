# charm-form

> Declarative form generator for terminal UIs -- define forms in YAML, get interactive CLI apps.

[Leia em Português](README.pt-br.md)

## Problem

Building interactive CLI forms requires writing Go code with specific TUI libraries. Every new form means more boilerplate, more compile cycles, and Go expertise just to ask a few questions.

## Solution

**charm-form** lets you define forms in plain YAML and get a fully interactive terminal UI. No Go knowledge required to create or modify forms -- just edit a YAML file.

The module ships with three pieces:

| Piece | Import path | What it does |
|---|---|---|
| Spec parser and validator | `github.com/junhinhow/charm-form` | Loads the YAML, validates it and applies the `validate` regex |
| Interactive renderer | `github.com/junhinhow/charm-form/renderer` | Bubble Tea v2 model that renders the fields and collects the answers |
| CLI | `github.com/junhinhow/charm-form/cmd/charm-form` | Runs a form and writes the answers as JSON |

## Installation

Install the CLI:

```bash
go install github.com/junhinhow/charm-form/cmd/charm-form@latest
```

Or add as a library:

```bash
go get github.com/junhinhow/charm-form
```

> Requires Go 1.25 or newer -- the floor imposed by `charm.land/lipgloss/v2`.

## Quick start

```bash
charm-form -file examples/deploy.yaml
```

The CLI opens the form in the terminal, collects the answers and prints them as JSON:

```json
{
  "environment": "production",
  "version": "v1.4.0",
  "dry_run": false,
  "services": ["api", "web"]
}
```

## YAML Spec Reference

```yaml
title: My Form             # Required - form title
description: Some info     # Optional - displayed below the title
fields:                    # Required - at least one field
  - name: field_id         # Required - unique identifier
    type: input            # Required - input|select|confirm|multiselect|text
    label: Display Label   # Field label shown to user
    placeholder: hint      # Placeholder text (input/text only)
    options:               # Required for select/multiselect
      - option1
      - option2
    required: true         # Whether the field is required
    validate: "^regex$"    # Regex validation pattern (optional)
```

### Field Types

| Type | Description | Answer type |
|------|-------------|-------------|
| `input` | Single-line text input | `string` |
| `text` | Multi-line text input | `string` |
| `select` | Single choice from options | `string` |
| `multiselect` | Multiple choices from options | `[]string` |
| `confirm` | Yes/no boolean | `bool` |

### Validation rules

- `required: true` -- the field cannot be left empty.
- `validate: "<regex>"` -- the answer must match the regular expression. The pattern is compiled while the spec is loaded, so an invalid regex is reported right away instead of failing silently.
- `name` must be unique across the form; duplicated names are rejected while the spec is loaded.
- `type` is case insensitive (`Input`, `INPUT` and `input` are the same type) and is normalized to lower case.

## CLI

```
Usage: charm-form [options] [form.yaml]

  -file string
        path to the YAML file with the form definition
  -out string
        file the result is written to (default: stdout)
  -json bool
        emit the answers as JSON; use -json=false for a readable summary
```

## Keyboard

| Key | Action |
|---|---|
| `enter` | Confirm the current field |
| `tab` / `shift+tab` | Next / previous field |
| `up` / `down` | Move the selection (`select` and `multiselect`) |
| `space` | Toggle an option (`multiselect`) |
| `y` / `n` | Toggle a boolean (`confirm`) |
| `ctrl+s` | Confirm a multi-line `text` field (`enter` inserts a line break) |
| `ctrl+c` / `esc` | Cancel the form |

## Usage as a library

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

### Working with the spec only

`charmform.Run(spec)` prints the specification itself as JSON -- useful to inspect or convert a form definition:

```go
spec, err := charmform.LoadFromFile("deploy.yaml")
if err != nil {
    log.Fatal(err)
}
if err := charmform.Run(spec); err != nil {
    log.Fatal(err)
}
```

`charmform.WriteJSON(w, spec)` does the same into any `io.Writer`.

## Examples

See the [examples/](examples/) directory for sample YAML form definitions.

## Roadmap

- [x] Interactive TUI rendering with Bubble Tea v2
- [x] CLI tool (`charm-form -file deploy.yaml`)
- [x] JSON output
- [x] Regex validation (`validate`)
- [ ] Output formats: env vars, TOML
- [ ] Conditional fields (show/hide based on other answers)
- [ ] Default values and environment variable interpolation
- [ ] Form composition (include/extend other YAML forms)
- [ ] Validation beyond regex (min/max, email, URL, etc.)

## License

MIT
