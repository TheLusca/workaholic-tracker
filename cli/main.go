// Comando horas — CLI do Workaholic Tracker (Caderno de Horas).
// Compartilha o mesmo schema de banco de horas (.json) do app web
// caderno-de-horas.html — ver caderno-de-horas-resumo.md.
package main

import (
	"fmt"
	"os"
)

const usage = `horas — rastreador de horas trabalhadas por projeto (linha de comando)

uso:
  horas start <projeto>                      inicia o timer num projeto
  horas stop                                 para o timer e registra a sessão
  horas add <projeto> --hoje <dur> [--nota "..."]
  horas add <projeto> --data YYYY-MM-DD --duracao <dur> [--nota "..."]
                                              registra horas manualmente
  horas hoje                                 total de hoje por projeto
  horas semana                               total da semana atual por projeto
  horas painel --ano                         heatmap dos últimos 12 meses
  horas status                               mostra se o timer está rodando

  <dur> aceita: 2h, 1h30, 45m, ou minutos puros (ex.: 90)

opções globais:
  --banco <caminho>   caminho do arquivo do banco de horas
                      (padrão: $HORAS_BANK ou ~/.local/share/horas/banco-de-horas.json)
`

func main() {
	args := os.Args[1:]

	var bankPath string
	filtered := args[:0:0]
	for i := 0; i < len(args); i++ {
		if args[i] == "--banco" {
			if i+1 >= len(args) {
				fail("--banco precisa de um caminho de arquivo")
			}
			bankPath = args[i+1]
			i++
			continue
		}
		filtered = append(filtered, args[i])
	}
	args = filtered

	if len(args) == 0 {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := args[0]
	rest := args[1:]

	if cmd == "-h" || cmd == "--help" || cmd == "help" {
		fmt.Print(usage)
		return
	}

	path, err := resolveBankPath(bankPath)
	if err != nil {
		fail("%v", err)
	}
	b, err := loadBank(path)
	if err != nil {
		fail("%v", err)
	}

	switch cmd {
	case "start":
		cmdStart(b, rest)
	case "stop":
		cmdStop(b, rest)
	case "add":
		cmdAdd(b, rest)
	case "hoje":
		cmdHoje(b, rest)
	case "semana":
		cmdSemana(b, rest)
	case "painel":
		cmdPainel(b, rest)
	case "status":
		cmdStatus(b)
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n\n", cmd)
		fmt.Print(usage)
		os.Exit(1)
	}
}
