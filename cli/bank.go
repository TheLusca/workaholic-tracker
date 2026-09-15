package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const schemaVersion = 1
const appID = "workaholic-tracker"

// mesma paleta usada em caderno-de-horas.html, pra projetos criados pela CLI
// ficarem com cor consistente se o banco for aberto depois no app web.
var palette = []string{"#D6923F", "#5C7D71", "#B5533C", "#7A8CAE", "#8C7A9E", "#A3A15C"}

// Project e Entry espelham exatamente o schema do banco de horas usado pelo
// app web (caderno-de-horas.html / caderno-de-horas-resumo.md), pra um mesmo
// arquivo .json poder ser aberto pelos dois.
type Project struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Entry struct {
	ID              string `json:"id"`
	ProjectID       string `json:"projectId"`
	Date            string `json:"date"` // YYYY-MM-DD
	DurationMinutes int    `json:"durationMinutes"`
	Note            string `json:"note"`
	Source          string `json:"source"` // "timer" | "manual"
}

type ActiveTimer struct {
	ProjectID string `json:"projectId"`
	StartedAt int64  `json:"startedAt"` // epoch ms, igual a Date.now() no JS
}

type Bank struct {
	SchemaVersion int          `json:"schemaVersion"`
	App           string       `json:"app"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
	Revision      int          `json:"revision"`
	Projects      []Project    `json:"projects"`
	Entries       []Entry      `json:"entries"`
	ActiveTimer   *ActiveTimer `json:"activeTimer"`

	path   string // caminho de onde foi carregado/onde será salvo
	loaded bool   // true se o arquivo já existia em disco
}

func nowISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

func genID() string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 36)
	suffix := strconv.FormatInt(int64(rand.Intn(60466176)), 36) // até 5 chars em base36
	return ts + suffix
}

// defaultBankPath segue a convenção documentada em caderno-de-horas-resumo.md:
// ~/.local/share/horas/banco-de-horas.json — pode ser sobrescrito com --banco
// ou a variável de ambiente HORAS_BANK.
func defaultBankPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("não foi possível localizar o diretório do usuário: %w", err)
	}
	return filepath.Join(home, ".local", "share", "horas", "banco-de-horas.json"), nil
}

func resolveBankPath(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if env := os.Getenv("HORAS_BANK"); env != "" {
		return env, nil
	}
	return defaultBankPath()
}

// loadBank lê o banco do disco. Se o arquivo não existir, retorna um banco
// vazio em memória (loaded=false) sem tocar o disco — só é gravado quando
// uma operação de fato muda algo e chama Save().
func loadBank(path string) (*Bank, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		now := nowISO()
		return &Bank{
			SchemaVersion: schemaVersion,
			App:           appID,
			CreatedAt:     now,
			UpdatedAt:     now,
			Revision:      0,
			Projects:      []Project{},
			Entries:       []Entry{},
			ActiveTimer:   nil,
			path:          path,
			loaded:        false,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao ler banco de horas em %s: %w", path, err)
	}
	var b Bank
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("banco de horas em %s não é um JSON válido: %w", path, err)
	}
	// leitura tolerante a arquivos sem metadados de schema (ex.: banco antigo)
	if b.SchemaVersion == 0 {
		b.SchemaVersion = schemaVersion
	}
	if b.App == "" {
		b.App = appID
	}
	now := nowISO()
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	if b.UpdatedAt == "" {
		b.UpdatedAt = now
	}
	if b.Projects == nil {
		b.Projects = []Project{}
	}
	if b.Entries == nil {
		b.Entries = []Entry{}
	}
	b.path = path
	b.loaded = true
	return &b, nil
}

// Save grava o banco no disco, atualizando updatedAt e incrementando revision
// — mesma semântica de saveBank()/writeFileHandle() no app web.
func (b *Bank) Save() error {
	if err := os.MkdirAll(filepath.Dir(b.path), 0o755); err != nil {
		return fmt.Errorf("erro ao criar diretório do banco de horas: %w", err)
	}
	b.UpdatedAt = nowISO()
	b.Revision++
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar banco de horas: %w", err)
	}
	tmp := b.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("erro ao gravar banco de horas: %w", err)
	}
	if err := os.Rename(tmp, b.path); err != nil {
		return fmt.Errorf("erro ao finalizar gravação do banco de horas: %w", err)
	}
	b.loaded = true
	return nil
}

func (b *Bank) findProject(name string) *Project {
	for i := range b.Projects {
		if b.Projects[i].Name == name {
			return &b.Projects[i]
		}
	}
	return nil
}

func (b *Bank) findProjectByID(id string) *Project {
	for i := range b.Projects {
		if b.Projects[i].ID == id {
			return &b.Projects[i]
		}
	}
	return nil
}

// findOrCreateProject casa por nome exato (igual ao import de CSV do app web);
// se não existir, cria um projeto novo com a próxima cor da paleta.
func (b *Bank) findOrCreateProject(name string) *Project {
	if p := b.findProject(name); p != nil {
		return p
	}
	p := Project{
		ID:    genID(),
		Name:  name,
		Color: palette[len(b.Projects)%len(palette)],
	}
	b.Projects = append(b.Projects, p)
	return &b.Projects[len(b.Projects)-1]
}
