// Comando charm-form executa um formulário declarativo definido em YAML.
//
// Uso:
//
//	charm-form -file formulario.yaml [-out respostas.json] [-json=false]
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	charmform "github.com/junhinhow/charm-form"
	"github.com/junhinhow/charm-form/renderer"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

// run executa o comando com os argumentos informados.
func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("charm-form", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		arquivo = fs.String("file", "", "caminho do arquivo YAML com a definição do formulário")
		saida   = fs.String("out", "", "arquivo onde gravar o resultado (padrão: stdout)")
		emJSON  = fs.Bool("json", true, "emitir as respostas em JSON; use -json=false para um resumo legível")
	)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Uso: charm-form [opções] [arquivo.yaml]\n\n")
		fmt.Fprintf(stderr, "Executa o formulário declarativo definido em YAML e emite as respostas.\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	caminho := *arquivo
	if caminho == "" && fs.NArg() > 0 {
		caminho = fs.Arg(0)
	}
	if caminho == "" {
		fs.Usage()
		return errors.New("informe o arquivo do formulário com -file")
	}

	spec, err := charmform.LoadFromFile(caminho)
	if err != nil {
		return err
	}

	respostas, err := renderer.Run(spec)
	if err != nil {
		if errors.Is(err, renderer.ErrCancelado) {
			return errors.New("formulário cancelado pelo usuário")
		}
		return err
	}

	destino := stdout
	if *saida != "" {
		arquivoSaida, err := os.Create(*saida)
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de saída %s: %w", *saida, err)
		}
		defer arquivoSaida.Close()
		destino = arquivoSaida
	}

	if *emJSON {
		if err := respostas.WriteJSON(destino); err != nil {
			return err
		}
	} else if err := escreveResumo(destino, spec, respostas); err != nil {
		return err
	}

	return nil
}

// escreveResumo grava as respostas em um formato legível, uma linha por
// campo, na ordem em que os campos aparecem na especificação.
func escreveResumo(w io.Writer, spec charmform.FormSpec, respostas charmform.Answers) error {
	for _, f := range spec.Fields {
		valor, ok := respostas[f.Name]
		if !ok {
			continue
		}
		texto := formatarValor(valor)
		if _, err := fmt.Fprintf(w, "%s: %s\n", f.Name, texto); err != nil {
			return fmt.Errorf("erro ao escrever resumo: %w", err)
		}
	}
	return nil
}

// formatarValor converte o valor de um campo em texto legível.
func formatarValor(valor any) string {
	switch v := valor.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		if v {
			return "sim"
		}
		return "não"
	case []string:
		return strings.Join(v, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}
