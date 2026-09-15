```
█   █  ███  ████  █   █  ███  █   █  ███  █     ███  ███     █████ ████   ███   ███  █   █ █████ ████    
█░  █░█ ░░█ █░░░█ █░ █ ░█ ░░█ █░  █░█ ░░█ █░     █░░█ ░░░     ░█░░░█░░░█ █ ░░█ █ ░░░ █░ █ ░█░░░░░█░░░█   
█░█ █░█░ ░█░████░░███ ░ █████░█████░█░ ░█░█░░    █░░█░ ░░░     █░░░████░░█████░█░ ░░░███ ░ ████░░████░░  
██░██░█░░ █░█░░█░ █░░█ ░█░░░█░█░░░█░█░░ █░█░░    █░░█░░        █░░ █░░█░ █░░░█░█░░   █░░█ ░█░░░░ █░░█░ ░ 
█░░ █░░███ ░█░░░█░█░░░█ █░░░█░█░░░█░░███ ░█████ ███░ ███       █░░ █░░░█░█░░░█░░███  █░░░█ █████░█░░░█░  
 ░░░ ░░ ░░░ ░░░  ░ ░░  ░ ░░  ░░░░  ░░ ░░░ ░░░░░░ ░░░  ░░░       ░░  ░░  ░ ░░  ░░ ░░░  ░░  ░ ░░░░░ ░░  ░  
  ░   ░  ░░░  ░   ░ ░   ░ ░   ░ ░   ░  ░░░  ░░░░░ ░░░  ░░░       ░   ░   ░ ░   ░  ░░░  ░   ░ ░░░░░ ░   ░ 
```

# Workaholic Tracker (Caderno de Horas)

Rastreador pessoal de horas trabalhadas por projeto, inspirado nos painéis de contribuição do GitHub/GitLab: controle **micro** (cada sessão registrada) + visualização **macro** (intensidade de trabalho ao longo do tempo).

## Como usar (app web)
1. Abra `caderno-de-horas.html` no Chrome, Edge ou Opera (usa File System Access API).
2. Clique em **"+ novo banco"** para criar seu banco de horas local (`.json`), ou **"carregar banco de horas"** para abrir um já existente.
3. Registre horas pelo timer ou por entrada manual; acompanhe pelo painel (dia / semana / mês / 6 meses / ano).

Sem esse banco carregado ou criado, nenhuma hora pode ser registrada — não há fallback em memória.

## Como usar (CLI)
```
cd cli
go build -o horas .
./horas start "Nome do Projeto"
./horas status
```
Lê e grava o **mesmo schema** de banco de horas do app web — dá pra apontar `--banco <caminho>` pro mesmo arquivo `.json`. Detalhes em [cli/README.md](cli/README.md).

## Stack
- App web: HTML/CSS/JS puro, sem build, sem backend. Único recurso externo: fontes via Google Fonts CDN.
- CLI: Go, binário único, sem runtime.

## Documentação
- [caderno-de-horas-resumo.md](caderno-de-horas-resumo.md) — visão geral do projeto, funcionalidades e arquitetura de armazenamento.
- [caderno-de-horas-design-system.md](caderno-de-horas-design-system.md) — paleta de cores, tipografia, ícones e demais recursos visuais.
- [cli/README.md](cli/README.md) — build, uso e comandos da CLI.

## Status
App web funcional para uso pessoal. CLI com primeira versão implementada (todos os comandos do design original), interoperando de verdade com o mesmo arquivo de banco de horas do app web.
