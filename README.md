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

## Como usar
1. Abra `caderno-de-horas.html` no Chrome, Edge ou Opera (usa File System Access API).
2. Clique em **"+ novo banco"** para criar seu banco de horas local (`.json`), ou **"carregar banco de horas"** para abrir um já existente.
3. Registre horas pelo timer ou por entrada manual; acompanhe pelo painel (dia / semana / mês / 6 meses / ano).

Sem esse banco carregado ou criado, nenhuma hora pode ser registrada — não há fallback em memória.

## Stack
HTML/CSS/JS puro, sem build, sem backend. Único recurso externo: fontes via Google Fonts CDN.

## Documentação
- [caderno-de-horas-resumo.md](caderno-de-horas-resumo.md) — visão geral do projeto, funcionalidades e arquitetura de armazenamento.
- [caderno-de-horas-design-system.md](caderno-de-horas-design-system.md) — paleta de cores, tipografia, ícones e demais recursos visuais.

## Status
Protótipo funcional para uso pessoal. Próximo passo: versão CLI (conceito já especificado no resumo do projeto).
