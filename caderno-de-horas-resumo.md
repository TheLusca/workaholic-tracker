# Caderno de Horas — resumo do projeto

## Ideia
Rastreador pessoal de horas trabalhadas por projeto, inspirado nos painéis de contribuição do GitHub/GitLab: controle **micro** (registro detalhado de cada sessão) + visualização **macro** (intensidade de trabalho ao longo do tempo).

## Identidade do projeto
- **Nome oficial (identificador técnico)**: `workaholic-tracker` — usado no campo `"app"` do schema do banco de horas e como nome do banco IndexedDB (`workaholic-tracker-fs`).
- **Nome exibido pro usuário**: "Caderno de Horas" (título da aba, cabeçalho) — continua em pt-BR, é a marca/nome amigável. Decisão consciente de manter os dois separados: identificador técnico em inglês/kebab-case, nome de exibição em português.
- **Identidade visual** (paleta de cores, tipografia, ícones/assets, espaçamento, animações): documentada em **[caderno-de-horas-design-system.md](caderno-de-horas-design-system.md)**.

## Status atual / próximos passos
- Protótipo web (`caderno-de-horas.html`) está funcionalmente completo: registro (timer + manual), painel dia/semana/mês/6 meses/ano, estatísticas, CSV (export + import), gerenciamento de projetos (criar/editar/excluir), banco de horas em arquivo `.json` selecionado diretamente (schema versionado) com trava/backup estilo Office implementada via IndexedDB (ver seção Armazenamento).
- **Correção de rumo**: a primeira versão da trava/backup pedia acesso à pasta inteira (`showDirectoryPicker`) pra poder criar um arquivo irmão oculto de verdade. Voltamos pra seleção direta do arquivo (sem pedir pasta) e reimplementamos a trava/backup como espelho no IndexedDB — mesma função, permissão bem menor.
- Testado com estado simulado e um `FileSystemFileHandle` simulado no navegador (bypassando os diálogos nativos, que não são automatizáveis pelas ferramentas usadas até aqui): todas as views do painel, edição/exclusão de projeto, e o ciclo completo de conectar/salvar/fechar/detectar-trava/recuperar-por-revisão do banco funcionam sem erros. **Falta validação manual num Chrome/Edge de verdade** só da interação real com os diálogos nativos de arquivo — a lógica por trás já foi validada.
- **Próxima implementação: a versão CLI** — já tem conceito desenhado (comandos, linguagem cogitada, referências de UX) na seção "Conceito de versão CLI" logo abaixo. Nada de código ainda, só a especificação.

## Decisões de plataforma
- MVP: web app client-side, sem backend, uso pessoal.
- Evolução futura cogitada: empacotar com **Tauri** para virar app desktop com ícone na bandeja do sistema (start/stop rápido do timer).
- Também desenhamos o conceito de uma versão **CLI** (não implementada, só especificada — ver abaixo).

## Modelo de dados
```
Project   { id, name, color }
TimeEntry { id, projectId, date (YYYY-MM-DD), durationMinutes, note, source: 'timer' | 'manual' }
```
Timer ativo é persistido à parte: `{ projectId, startedAt }` — permite retomar a contagem se a página recarregar.

## Funcionalidades do protótipo (já implementadas em `caderno-de-horas.html`)
1. **Duas formas de registrar horas, em abas de verdade**: ⏱ Timer (start/stop) e ✎ Entrada manual (data + projeto + horas/min + nota) — só uma fica visível por vez.
2. Indicador visual: um pontinho vermelho pulsante aparece ao lado da aba "⏱ timer" sempre que o timer está rodando, mesmo com a aba "✎ entrada" selecionada — assim dá pra notar o timer ativo em segundo plano sem trocar de aba.
3. **Painel de visualização com 5 granularidades alternáveis**: Dia / Semana / Mês / 6 meses / Ano.
   - **Dia**: barras horizontais por projeto no dia selecionado + total, com navegação ◀ ▶.
   - **Semana**: gráfico de barras empilhadas por dia (seg–dom), coloridas por projeto, com legenda.
   - **Mês**: calendário do mês (grade seg–dom, células coloridas por intensidade igual ao heatmap), dias fora do mês esmaecidos, dia atual com contorno, navegação ◀ ▶ por mês, tooltip por dia igual ao heatmap, total do mês.
   - **6 meses / Ano**: heatmap estilo GitHub (grade semanas × dias da semana); tooltip ao passar o mouse mostra a divisão por projeto naquele dia, não só o total.
   - Filtro por projeto ou "todos os projetos" (afeta semana/mês/6 meses/ano).
4. Estatísticas: total registrado, sequência de dias seguidos (streak), horas da semana atual.
5. Lista de registros recentes (até 40), com edição e exclusão inline.
6. **Exportação e importação de planilha CSV** — export com separador `;`, decimais com vírgula, BOM UTF-8 (acentos corretos), pensado para abrir sem bagunça no Excel/Google Sheets em pt-BR. Import é o reverso: lê esse mesmo formato e recria os registros no banco atual.
   - Projetos são casados por **nome exato**; se o nome não existir no banco atual, um projeto novo é criado automaticamente (cor da paleta por ordem). Linhas com projeto vazio ou `"projeto removido"` (placeholder de projeto excluído) são ignoradas — não há como reconstruir a referência original.
   - Duração é lida da coluna "Duracao (min)" (inteira, fonte de verdade); só cai pra coluna de horas decimais se a de minutos estiver ausente/inválida.
   - Como o CSV não carrega os `id`s originais, todo registro importado ganha um `id` novo — reimportar o mesmo CSV duas vezes duplica os registros (nenhuma deduplicação é feita). Antes de aplicar, mostra um resumo (quantos registros, quantos projetos novos, quantas linhas ignoradas) pra confirmação.
7. **Gerenciamento de projetos** via popover "projetos" no cabeçalho: lista todos os projetos existentes, cada um com botões "editar" (edição inline: nome + escolha de cor entre as 6 da paleta fixa, com salvar/cancelar) e "excluir" (com confirmação; bloqueado se o timer estiver rodando nesse projeto), e um campo fixo no rodapé do popover pra criar projeto novo (cor atribuída automaticamente por ordem). Desabilitado enquanto nenhum banco de horas estiver carregado/criado (ver Armazenamento).
   - **Excluir projeto não apaga o histórico**: os `TimeEntry` que referenciavam aquele projeto continuam no banco, só passam a aparecer como "projeto removido" (cinza) em todas as views — mesmo comportamento gracioso que já existia para o caso de dado órfão.
   - **Decisão consciente sobre cor**: a cor é atributo do **projeto**, não do lançamento (`TimeEntry` não guarda cor própria). Isso significa que editar a cor de um projeto reexibe *todo* o histórico dele com a cor nova — não há "cor congelada por lançamento". Avaliado e considerado aceitável: faz mais sentido cor pertencer ao projeto (identidade visual consistente) do que ao registro individual.

## Direção visual
Conceito: "painel de instrumentos" — paleta escura e quente com acentos em latão/âmbar. Referência completa (todas as cores, tipografia, ícones/assets, espaçamento, animações, responsividade) foi extraída pra um arquivo dedicado: **[caderno-de-horas-design-system.md](caderno-de-horas-design-system.md)**.

## Armazenamento — "banco de horas"
Não existe fallback em memória nem em `window.storage` (Claude.ai): os dados só existem dentro de um arquivo `.json` local ("banco de horas") escolhido pelo usuário, e **enquanto nenhum banco estiver carregado ou criado, nenhuma hora pode ser registrada** — a interface de registro/painel/estatísticas/registros fica escondida por trás de um painel de bloqueio, e o botão "projetos" fica desabilitado.

### Schema do arquivo `.json`
```json
{
  "schemaVersion": 1,
  "app": "workaholic-tracker",
  "createdAt": "2026-08-08T18:00:00.000Z",
  "updatedAt": "2026-08-08T19:30:00.000Z",
  "revision": 42,
  "projects": [],
  "entries": [],
  "activeTimer": null
}
```
- `createdAt` é gravado uma vez na criação do banco e nunca muda depois.
- `updatedAt` e `revision` são atualizados a cada `saveBank()` (uma escrita = uma revisão a mais).
- Ao abrir um banco antigo sem esses campos de metadado, o app assume valores default (`schemaVersion: 1`, `revision: 0`, `createdAt`/`updatedAt` = agora) em vez de falhar — leitura tolerante a formatos anteriores.

### Conexão: seleção direta do arquivo
Implementado via File System Access API do navegador, selecionando o arquivo `.json` diretamente (`showOpenFilePicker`/`showSaveFilePicker`) — **não pede acesso a pasta**. Chegamos a implementar uma versão baseada em `showDirectoryPicker` (pra poder criar um arquivo irmão oculto de verdade na pasta, ver histórico abaixo), mas voltamos atrás: pedir a pasta inteira só pra viabilizar a trava/backup trocava uma permissão pequena (um arquivo) por uma grande (pasta inteira), e a expectativa era continuar selecionando o arquivo diretamente.
- **"carregar banco de horas"** — `showOpenFilePicker` (fluxo de leitura) pra escolher um `.json` **já existente**. Nunca cria nem sobrescreve nada nesse passo.
- **"+ novo banco"** — `showSaveFilePicker` (fluxo de criação), só grava vazio se o arquivo escolhido realmente não tiver conteúdo. Se o usuário escolher um arquivo que já tem dados, o código carrega o conteúdo em vez de sobrescrever.
- Depois de conectado, toda mutação (`saveBank()`) grava o objeto completo (schema acima) no arquivo, e o handle fica guardado num IndexedDB próprio (`workaholic-tracker-fs`) pra tentar reconectar sozinho na próxima abertura.
- **Isso só reconecta sem clique se o navegador ainda considerar a permissão "granted"** — em `file://` (e às vezes entre reloads) o Chrome costuma esquecer essa permissão; nesse caso o app mostra "🔌 reconectar [nome do arquivo]" em vez de voltar pro estado bloqueado sem explicação — um clique renova a permissão sem reabrir o seletor.
- Botão do cabeçalho vira `📁 nome-do-arquivo.json` quando conectado; clicar nele fecha o banco (com confirmação) e volta pro estado bloqueado — os dados continuam salvos no arquivo, só o app "esquece" dele.
- **Só funciona em Chrome/Edge/Opera** (File System Access API não existe no Firefox/Safari). Nesses navegadores o app mostra um aviso no painel de bloqueio e fica bloqueado permanentemente — decisão consciente, não um bug: preferimos travar com aviso claro a arriscar perda de dados com fallback silencioso em memória.

### Cópia oculta estilo Microsoft Office — trava + backup (via IndexedDB, não arquivo de verdade)
Como a conexão é por arquivo solto (não pasta), não há como criar um arquivo irmão real ao lado do banco — a API não expõe "pasta pai" a partir de um handle de arquivo, por design de segurança. A cópia oculta é implementada como um espelho no **IndexedDB do navegador** (`writeShadowCopy()`), guardando `{ fileHandle, lock: { app, lockedAt }, snapshot: <mesmo schema do arquivo principal> }`. Cumpre as mesmas duas funções do `~$arquivo.docx` do Word, só que sem ser um arquivo visível no Explorer:
1. **Trava**: ao conectar num arquivo, se o IndexedDB já tiver uma trava registrada pra esse mesmo arquivo (`FileSystemFileHandle.isSameEntry`), o app avisa que o banco pode já estar aberto em outra aba/janela deste navegador, ou que a sessão anterior não fechou corretamente — cita o horário da trava e pergunta se quer continuar mesmo assim.
2. **Recuperação por revisão**: se a cópia oculta tiver uma `revision` **maior** que a do arquivo recém-lido (sinal de que a última gravação no arquivo pode ter falhado), o app oferece recuperar os dados a partir da cópia local antes de prosseguir.
- Atualizada junto com todo `saveBank()` (mesmo `revision`/conteúdo do arquivo principal) e ao conectar/reconectar; **limpa automaticamente** só quando o banco é fechado corretamente pelo botão do cabeçalho (`closeBank()`) — se o app fechar de outro jeito (fechar aba, crash), a trava fica pra trás de propósito, é o sinal de alerta pro próximo "abrir".
- Limitação importante: como é uma trava por navegador (IndexedDB não é compartilhado entre navegadores/dispositivos), ela não detecta o arquivo sendo editado por *outro navegador* ou *outra máquina* — só a mesma instalação de Chrome/Edge que já abriu esse arquivo antes.
- Testado com um `FileSystemFileHandle` simulado em memória (já que o diálogo nativo não é automatizável): conexão cria o schema + a cópia oculta corretamente, a recuperação por revisão mais alta funciona (testado forçando o arquivo principal a regredir pra uma revisão antiga), a detecção de trava dispara e interrompe a conexão se o usuário cancelar, e o fechamento limpo remove a trava do IndexedDB.

### Caminhos de armazenamento ainda em aberto
- CLI: SQLite ou JSON em `~/.local/share/horas/`.
- Se a versão desktop evoluir com Tauri, o arquivo `.json` via navegador deixa de ser necessário — o próprio Tauri dá acesso a filesystem sem as limitações do navegador (Chromium-only, permissão por sessão), inclusive pra criar arquivos de trava reais ao lado do banco se fizer sentido nessa versão.

## Conceito de versão CLI (desenhada, não implementada)
```
horas start <projeto>
horas stop
horas add <projeto> --hoje 2h --nota "..."
horas hoje
horas semana
horas painel --ano
horas status
```
- Linguagem cogitada: Go ou Rust (binário único, sem runtime) vs Node/Python (prototipagem mais rápida).
- Referências de UX citadas: `timewarrior`, WakaTime CLI.
- Ponto de atenção de design: um comando `status` sempre rápido de consultar, pra não esquecer o timer rodando.

## Arquivo do protótipo
`caderno-de-horas.html` — HTML/CSS/JS puro, arquivo único, sem dependências externas de build (só carrega fontes do Google Fonts via CDN). Serve de referência funcional de UI/UX e lógica de negócio (cálculo de streak, agregação por dia/semana, heatmap) ao portar para outra stack.
