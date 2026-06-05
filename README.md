# charm-form

> Declarative form generator for terminal UIs -- define forms in YAML, get interactive CLI apps.

[Leia em Portugues](README.pt-br.md)

## Problem

Building interactive CLI forms requires writing Go code with specific TUI libraries. Every new form means more boilerplate, more compile cycles, and Go expertise just to ask a few questions.

## Solution

**charm-form** lets you define forms in plain YAML and get a fully interactive terminal UI. No Go knowledge required to create or modify forms -- just edit a YAML file.

## Installation

```bash
go install github.com/junhinhow/charm-form@latest
```

Or add as a library:

```bash
go get github.com/junhinhow/charm-form
```

## YAML Spec Reference

```yaml
title: My Form             # Required - form title
description: Some info     # Optional - displayed below title
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

| Type | Description |
|------|-------------|
| `input` | Single-line text input |
| `text` | Multi-line text input |
| `select` | Single choice from options |
| `multiselect` | Multiple choices from options |
| `confirm` | Yes/no boolean |

## Usage

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

## Examples

See the [examples/](examples/) directory for sample YAML form definitions.

## Roadmap

- [ ] Interactive TUI rendering with [charmbracelet/huh](https://github.com/charmbracelet/huh)
- [ ] CLI tool (`charm-form run deploy.yaml`)
- [ ] Output formats: JSON, env vars, TOML
- [ ] Conditional fields (show/hide based on other answers)
- [ ] Default values and environment variable interpolation
- [ ] Form composition (include/extend other YAML forms)
- [ ] Validation beyond regex (min/max, email, URL, etc.)

## License

MIT
