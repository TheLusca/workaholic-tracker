package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRe = regexp.MustCompile(`^(?:(\d+(?:\.\d+)?)h)?(?:(\d+)m?)?$`)

// parseDuration aceita "2h", "1h30", "1h30m", "45m" ou minutos puros ("90").
func parseDuration(s string) (int, error) {
	s = strings.TrimSpace(s)
	m := durationRe.FindStringSubmatch(s)
	if s == "" || m == nil || (m[1] == "" && m[2] == "") {
		return 0, fmt.Errorf("duração inválida: %q (use algo como 2h, 1h30, 45m ou minutos puros)", s)
	}
	minutes := 0.0
	if m[1] != "" {
		h, _ := strconv.ParseFloat(m[1], 64)
		minutes += h * 60
	}
	if m[2] != "" {
		mm, _ := strconv.Atoi(m[2])
		minutes += float64(mm)
	}
	total := int(minutes + 0.5)
	if total <= 0 {
		return 0, fmt.Errorf("duração precisa ser maior que zero: %q", s)
	}
	return total, nil
}

// formatDuration espelha formatDuration() do app web (ex.: "1h30", "45min", "2h").
func formatDuration(mins int) string {
	h := mins / 60
	m := mins % 60
	switch {
	case h == 0:
		return fmt.Sprintf("%dmin", m)
	case m == 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dh%02d", h, m)
	}
}

func todayDate() string {
	return time.Now().Format("2006-01-02")
}

func dateFromEpochMs(ms int64) string {
	return time.UnixMilli(ms).Format("2006-01-02")
}

var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func validateDate(s string) (string, error) {
	if !dateRe.MatchString(s) {
		return "", fmt.Errorf("data inválida: %q (use o formato YYYY-MM-DD)", s)
	}
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return "", fmt.Errorf("data inválida: %q", s)
	}
	return s, nil
}

// weekStart retorna a segunda-feira da semana de d, à meia-noite.
func weekStart(d time.Time) time.Time {
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	dow := int(d.Weekday()) // domingo=0
	diff := 1 - dow
	if dow == 0 {
		diff = -6
	}
	return d.AddDate(0, 0, diff)
}

func formatDatePtFull(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02/01/2006")
}
