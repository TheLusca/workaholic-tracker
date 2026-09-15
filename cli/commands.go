package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "erro: "+format+"\n", args...)
	os.Exit(1)
}

// ---------- horas start <projeto> ----------
func cmdStart(b *Bank, args []string) {
	if len(args) < 1 {
		fail("uso: horas start <projeto>")
	}
	name := args[0]
	if b.ActiveTimer != nil {
		proj := b.findProjectByID(b.ActiveTimer.ProjectID)
		started := time.UnixMilli(b.ActiveTimer.StartedAt).Format("15:04")
		pname := "?"
		if proj != nil {
			pname = proj.Name
		}
		fail("timer já está rodando em %q desde %s. use 'horas stop' primeiro.", pname, started)
	}
	proj := b.findOrCreateProject(name)
	b.ActiveTimer = &ActiveTimer{ProjectID: proj.ID, StartedAt: time.Now().UnixMilli()}
	if err := b.Save(); err != nil {
		fail("%v", err)
	}
	fmt.Printf("⏱  timer iniciado em %q\n", proj.Name)
}

// ---------- horas stop ----------
func cmdStop(b *Bank, args []string) {
	if b.ActiveTimer == nil {
		fail("nenhum timer rodando.")
	}
	elapsedMs := time.Now().UnixMilli() - b.ActiveTimer.StartedAt
	minutes := int(elapsedMs / 60000)
	if minutes < 1 {
		minutes = 1
	}
	p := b.findProjectByID(b.ActiveTimer.ProjectID)
	pname := "projeto removido"
	if p != nil {
		pname = p.Name
	}
	entry := Entry{
		ID:              genID(),
		ProjectID:       b.ActiveTimer.ProjectID,
		Date:            dateFromEpochMs(b.ActiveTimer.StartedAt),
		DurationMinutes: minutes,
		Note:            "",
		Source:          "timer",
	}
	b.Entries = append(b.Entries, entry)
	b.ActiveTimer = nil
	if err := b.Save(); err != nil {
		fail("%v", err)
	}
	fmt.Printf("■  parado: %q — %s\n", pname, formatDuration(minutes))
}

// ---------- horas add <projeto> --hoje <duração> | --data YYYY-MM-DD --duracao <duração> [--nota "..."] ----------
func cmdAdd(b *Bank, args []string) {
	if len(args) < 1 || len(args[0]) == 0 || args[0][0] == '-' {
		fail("uso: horas add <projeto> --hoje 2h [--nota \"...\"]  (ou --data YYYY-MM-DD --duracao 2h)")
	}
	name := args[0]
	rest := args[1:]

	var hojeVal, dataVal, duracaoVal, nota string
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--hoje":
			i++
			if i >= len(rest) {
				fail("--hoje precisa de um valor de duração, ex.: --hoje 2h")
			}
			hojeVal = rest[i]
		case "--data":
			i++
			if i >= len(rest) {
				fail("--data precisa de uma data, ex.: --data 2026-09-10")
			}
			dataVal = rest[i]
		case "--duracao":
			i++
			if i >= len(rest) {
				fail("--duracao precisa de um valor, ex.: --duracao 2h")
			}
			duracaoVal = rest[i]
		case "--nota":
			i++
			if i >= len(rest) {
				fail("--nota precisa de um texto")
			}
			nota = rest[i]
		default:
			fail("opção desconhecida: %s", rest[i])
		}
	}

	var date string
	var durStr string
	switch {
	case hojeVal != "":
		date = todayDate()
		durStr = hojeVal
	case dataVal != "" && duracaoVal != "":
		var err error
		date, err = validateDate(dataVal)
		if err != nil {
			fail("%v", err)
		}
		durStr = duracaoVal
	default:
		fail("informe --hoje <duração> ou --data YYYY-MM-DD --duracao <duração>")
	}

	minutes, err := parseDuration(durStr)
	if err != nil {
		fail("%v", err)
	}

	proj := b.findOrCreateProject(name)
	entry := Entry{
		ID:              genID(),
		ProjectID:       proj.ID,
		Date:            date,
		DurationMinutes: minutes,
		Note:            nota,
		Source:          "manual",
	}
	b.Entries = append(b.Entries, entry)
	if err := b.Save(); err != nil {
		fail("%v", err)
	}
	fmt.Printf("✎  registrado: %q em %s — %s\n", proj.Name, formatDatePtFull(date), formatDuration(minutes))
}

// ---------- agregação compartilhada ----------
type projectTotal struct {
	name string
	mins int
}

func aggregateByProject(b *Bank, from, to string) []projectTotal {
	totals := map[string]int{}
	order := []string{}
	for _, e := range b.Entries {
		if e.Date < from || e.Date > to {
			continue
		}
		p := b.findProjectByID(e.ProjectID)
		name := "projeto removido"
		if p != nil {
			name = p.Name
		}
		if _, ok := totals[name]; !ok {
			order = append(order, name)
		}
		totals[name] += e.DurationMinutes
	}
	result := make([]projectTotal, 0, len(order))
	for _, name := range order {
		result = append(result, projectTotal{name: name, mins: totals[name]})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].mins > result[j].mins })
	return result
}

func printTotals(title string, totals []projectTotal) {
	fmt.Println(title)
	if len(totals) == 0 {
		fmt.Println("  nenhum registro.")
		return
	}
	total := 0
	for _, t := range totals {
		fmt.Printf("  %-30s %s\n", t.name, formatDuration(t.mins))
		total += t.mins
	}
	fmt.Printf("  %-30s %s\n", "total", formatDuration(total))
}

// ---------- horas hoje ----------
func cmdHoje(b *Bank, args []string) {
	today := todayDate()
	totals := aggregateByProject(b, today, today)
	printTotals(fmt.Sprintf("hoje (%s):", formatDatePtFull(today)), totals)
}

// ---------- horas semana ----------
func cmdSemana(b *Bank, args []string) {
	monday := weekStart(time.Now())
	sunday := monday.AddDate(0, 0, 6)
	from := monday.Format("2006-01-02")
	to := sunday.Format("2006-01-02")
	totals := aggregateByProject(b, from, to)
	printTotals(fmt.Sprintf("semana (%s – %s):", formatDatePtFull(from), formatDatePtFull(to)), totals)
}

// ---------- horas status ----------
func cmdStatus(b *Bank) {
	if b.ActiveTimer == nil {
		fmt.Println("○ nenhum timer rodando.")
		return
	}
	p := b.findProjectByID(b.ActiveTimer.ProjectID)
	name := "projeto removido"
	if p != nil {
		name = p.Name
	}
	elapsed := time.Since(time.UnixMilli(b.ActiveTimer.StartedAt))
	h := int(elapsed.Hours())
	m := int(elapsed.Minutes()) % 60
	s := int(elapsed.Seconds()) % 60
	fmt.Printf("● timer rodando em %q — %02d:%02d:%02d\n", name, h, m, s)
}

// ---------- horas painel --ano ----------
// caracteres ASCII puros (largura garantidamente única em qualquer terminal) —
// os blocos de sombreamento Unicode (·░▒▓█) têm largura "ambígua" e muitos
// terminais os renderizam em largura dupla, quebrando o alinhamento com os
// rótulos de mês (que são ASCII).
var levelChars = []string{".", ":", "+", "*", "#"}

func levelFor(mins int) int {
	switch {
	case mins == 0:
		return 0
	case mins < 60:
		return 1
	case mins < 180:
		return 2
	case mins < 360:
		return 3
	default:
		return 4
	}
}

func cmdPainel(b *Bank, args []string) {
	ano := false
	for _, a := range args {
		if a == "--ano" {
			ano = true
		}
	}
	if !ano {
		fail("uso: horas painel --ano  (por enquanto só a visão anual está implementada)")
	}

	minByDate := map[string]int{}
	for _, e := range b.Entries {
		minByDate[e.Date] += e.DurationMinutes
	}

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	start := today.AddDate(0, 0, -370)
	start = weekStart(start)

	// monta as semanas (colunas) com 7 dias cada (linhas: seg..dom)
	type cell struct {
		date string
		mins int
	}
	var days []cell
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		days = append(days, cell{date: key, mins: minByDate[key]})
	}
	weeks := len(days) / 7
	if len(days)%7 != 0 {
		weeks++
	}

	// rótulos de mês: um buffer com exatamente `weeks` colunas, a mesma
	// largura da grade abaixo, pra garantir alinhamento perfeito.
	labelRow := make([]byte, weeks)
	for i := range labelRow {
		labelRow[i] = ' '
	}
	lastMonth := time.Month(0)
	monthsPt := []string{"jan", "fev", "mar", "abr", "mai", "jun", "jul", "ago", "set", "out", "nov", "dez"}
	for w := 0; w < weeks; w++ {
		idx := w * 7
		if idx >= len(days) {
			break
		}
		d, _ := time.Parse("2006-01-02", days[idx].date)
		if d.Month() != lastMonth {
			label := monthsPt[int(d.Month())-1] // sempre ASCII puro (3 letras)
			for j := 0; j < len(label) && w+j < weeks; j++ {
				labelRow[w+j] = label[j]
			}
			lastMonth = d.Month()
		}
	}
	fmt.Println("    " + string(labelRow))

	weekdaysPt := []string{"seg", "ter", "qua", "qui", "sex", "sáb", "dom"}
	for row := 0; row < 7; row++ {
		fmt.Printf("%s ", weekdaysPt[row])
		for w := 0; w < weeks; w++ {
			idx := w*7 + row
			if idx >= len(days) {
				fmt.Print(" ")
				continue
			}
			lvl := levelFor(days[idx].mins)
			fmt.Print(levelChars[lvl])
		}
		fmt.Println()
	}
	fmt.Println("menos " + levelChars[0] + levelChars[1] + levelChars[2] + levelChars[3] + levelChars[4] + " mais")

	totals := aggregateByProject(b, days[0].date, days[len(days)-1].date)
	total := 0
	for _, t := range totals {
		total += t.mins
	}
	fmt.Printf("\ntotal no período: %s\n", formatDuration(total))
}
