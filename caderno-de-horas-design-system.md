# Caderno de Horas (workaholic-tracker) — design system

Referência visual do protótipo `caderno-de-horas.html`. Todo valor aqui foi extraído direto do CSS/HTML do arquivo — se algo mudar no código, atualize aqui também.

## Conceito
"Painel de instrumentos" — paleta escura e quente (grafite/carvão, não preto puro) com acentos em latão/âmbar, remetendo a um painel de controle/medidor de horas trabalhadas. Nada de imagens externas: todo recurso visual é CSS, SVG inline ou caractere Unicode/emoji.

## Cores

Definidas como custom properties em `:root`, todas em `caderno-de-horas.html`.

| Variável | Hex | Uso |
|---|---|---|
| `--bg` | `#1E1B18` | fundo da página |
| `--panel` | `#26221E` | fundo dos cards/painéis (`.panel`, `.stat`) |
| `--panel-2` | `#2E2924` | fundo de elementos "dentro" de um painel (inputs, selects, trilhos de barra, view-switch) |
| `--ink` | `#EDE6D6` | texto principal |
| `--muted` | `#8C8378` | texto secundário/label (uppercase, datas, placeholders) |
| `--brass` | `#D6923F` | cor de destaque primária (botões primários, aba ativa, seleção de texto) |
| `--brass-bright` | `#F2B565` | hover/estado ativo dos elementos em `--brass`, foco (`:focus-visible`), valores numéricos de destaque |
| `--teal` | `#5C7D71` | cor secundária (pouco usada diretamente na UI; base pra paleta de projetos) |
| `--rust` | `#B5533C` | estado "timer rodando" (botão vermelho pulsante), excluir (hover), indicador ao vivo |
| `--border` | `#3A342C` | bordas de painéis, inputs, divisores |
| `--lvl0` … `--lvl4` | `#2A2620` `#4A3B24` `#7A5A2C` `#B98536` `#F2B565` | escala de 5 níveis do heatmap/calendário mensal, do "sem atividade" ao "mais intenso" — tons de âmbar, não o verde do GitHub |

Cores fora do sistema de variáveis, usadas pontualmente:
- `#211D18` — texto sobre fundo `--brass`/`--brass-bright` (botão primário, aba ativa)
- `#F3E5E0` — texto do botão do timer quando "rodando" (sobre `--rust`)
- `#0F0D0B` — fundo do tooltip
- `rgba(0,0,0,.4)` — sombra do popover
- `rgba(181,83,60,.4→0)` — glow do `@keyframes pulse` (timer rodando)
- `rgba(242,181,101,.5)` — glow das células de heatmap/calendário no nível máximo (`lvl4`)
- `rgba(33,29,24,.7)` — cor do número do dia no calendário mensal quando a célula está clara (níveis 2–4), pra manter contraste

### Paleta de projetos
Array fixo usado pra atribuir cor automaticamente a projetos novos (`PALETTE` no JS, ciclo por índice — `state.projects.length % PALETTE.length`):
```
#D6923F  #5C7D71  #B5533C  #7A8CAE  #8C7A9E  #A3A15C
```
(latão, verde-azulado, terracota, azul acinzentado, roxo acinzentado, oliva — 6 tons dessaturados que convivem bem com o fundo escuro). É a mesma paleta oferecida como swatches na edição manual de cor do projeto.

## Tipografia
- **Space Grotesk** (400/500/600/700) — interface, títulos, rótulos, texto geral.
- **JetBrains Mono** (400/500/600/700) — números: timer, datas, durações, valores de estatística, revisão do banco. Classe utilitária `.mono` aplica em qualquer elemento.
- Carregadas via Google Fonts CDN: `https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600;700&display=swap` (único recurso externo do projeto além do próprio HTML).
- Fallbacks: `ui-sans-serif, system-ui, sans-serif` pra Space Grotesk; `ui-monospace, 'SF Mono', Consolas, monospace` pra JetBrains Mono.

## Ícones e símbolos
Não há biblioteca de ícones nem imagens/PNG/SVG externas — só:
- **Ícone de marca** (cabeçalho, ao lado de "Caderno de Horas"): SVG inline 24×24 (renderizado a 26×26px), um relógio estilizado — círculo (`stroke` `--brass`) com dois ponteiros em ângulo (`stroke` `--brass-bright`) e um ponto central. Não há favicon de arquivo separado.
- **Símbolos de texto/Unicode** usados como ícone funcional (sem `<svg>`, só caractere): `►` iniciar timer, `■` parar timer, `◀` `▶` navegação de painel, `⏱` aba timer, `✎` aba entrada manual/edição, `📁` banco conectado, `🔌` reconectar banco, `↑` importar CSV, `↓` exportar CSV.
- **Pontos coloridos** (`<span>` com `border-radius:50%` e `background` inline): identificam projeto por cor em várias listas (registros, dia, semana, popover de projetos) — não são ícones fixos, a cor vem do próprio projeto.

## Espaçamento, raio e sombra
- `--radius: 6px` — padrão de botões, inputs, selects, popover.
- Painéis (`.panel`) e cards de estatística (`.stat`): `border-radius: 10px`.
- Espaçamento entre seções do app: `gap: 16px` (`#app`, `.workspace`).
- Espaçamento entre controles dentro de um grupo: `8–10px`.
- Sombra do popover: `0 8px 24px rgba(0,0,0,.4)`.
- Glow (sombra colorida) em dois lugares: célula de heatmap/calendário no nível máximo, e no `save-indicator` ao "piscar" depois de salvar (`box-shadow: 0 0 8px var(--brass-bright)`).

## Animações
- **`pulse`** (1.8s, loop) — no botão do timer quando está rodando: anel de sombra em `--rust` que expande e desaparece.
- **`pulse-dot`** (1.4s, loop) — no ponto vermelho ao lado da aba "⏱ timer" quando o timer está ativo em segundo plano: opacidade entre 1 e 0.25.
- **Entrada dos painéis**: fade-in + leve translateY ao carregar a página, com atraso escalonado por painel (`.02s`, `.08s`, `.14s`, `.20s`) pra dar sensação de "montagem" sequencial.
- Tudo respeita `prefers-reduced-motion: reduce` (desliga toda animação/transição globalmente nesse caso).

## Responsividade
Breakpoint único em `max-width: 640px`: formulário manual vira 2 colunas, timer diminui a fonte (40px → 32px), stats empilham verticalmente, filtro de projeto ocupa a largura toda, grade do calendário mensal diminui a altura das células (36px → 28px).

## Assets externos
Único recurso carregado de fora: a folha de fontes do Google Fonts (link acima). Tudo mais — ícone de marca, cores, ilustrações — é gerado por CSS/SVG inline dentro do próprio `caderno-de-horas.html`. Isso é intencional: mantém o arquivo único e sem dependência de build, como documentado em [caderno-de-horas-resumo.md](caderno-de-horas-resumo.md).
