# horas — CLI do Workaholic Tracker

Implementação em Go da versão de linha de comando do Caderno de Horas. Lê e grava o **mesmo schema** de banco de horas (`.json`) usado pelo app web (`caderno-de-horas.html`) — dá pra apontar a CLI pro mesmo arquivo que você já usa no navegador.

## Build

```
cd cli
go build -o horas .
```

Gera um binário único (`horas` / `horas.exe`), sem runtime necessário.

## Banco de horas

Por padrão, o arquivo fica em `~/.local/share/horas/banco-de-horas.json` (criado automaticamente no primeiro registro). Pra usar outro caminho — por exemplo, o mesmo arquivo que você já abre no app web:

```
horas --banco "C:\caminho\para\banco-de-horas.json" status
```

ou defina a variável de ambiente `HORAS_BANK` pra não precisar repetir a flag.

## Comandos

```
horas start <projeto>                      inicia o timer num projeto
horas stop                                 para o timer e registra a sessão
horas add <projeto> --hoje <dur> [--nota "..."]
horas add <projeto> --data YYYY-MM-DD --duracao <dur> [--nota "..."]
horas hoje                                 total de hoje por projeto
horas semana                               total da semana atual por projeto
horas painel --ano                         heatmap ASCII dos últimos 12 meses
horas status                               mostra se o timer está rodando
```

`<dur>` aceita `2h`, `1h30`, `45m` ou minutos puros (`90`).

Projetos são criados automaticamente pelo nome na primeira vez que aparecem em `start`/`add` — não existe comando separado de "criar projeto" (mesma filosofia leve do design original).

## Status
Implementação inicial — todos os comandos da especificação original (`caderno-de-horas-resumo.md`, seção "Versão CLI") mais duas extensões pequenas: `--data`/`--duracao` em `add` (pra registrar dias que não são hoje) e `--banco`/`HORAS_BANK` (pra apontar pra um arquivo específico).

Testado manualmente: `go vet`/`gofmt` limpos; start/stop/status/add/hoje/semana/painel com um banco isolado; tratamento de erro (duração inválida, sem projeto, `stop` sem timer); interoperabilidade real com o `banco-de-horas.json` gerado pelo app web (leu os dados, reconheceu o projeto existente sem duplicar, preservou `createdAt` e incrementou `revision` corretamente).
