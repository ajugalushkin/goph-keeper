package version

import (
	// init embed package
	_ "embed"
	"io"
	"log/slog"
	"strings"
	"text/template"
)

//go:embed version.txt
var greeting string

var (
	// BuildVersion версия исполняемого файла
	BuildVersion string
	// BuildDate время сборки
	BuildDate string
	// BuildCommit коммит хеш
	BuildCommit string
)

type buildInfo struct {
	Version string
	Date    string
	Commit  string
}

// WriteBuildInfo записывает информацию о сборке
func WriteBuildInfo(w io.Writer) {
	info := buildInfo{
		Version: "N/A",
		Date:    "N/A",
		Commit:  "N/A",
	}

	if BuildVersion != "" {
		info.Version = BuildVersion
	}
	if BuildDate != "" {
		info.Date = BuildDate
	}
	if BuildCommit != "" {
		info.Commit = BuildCommit
	}

	tmpl := template.Must(template.New("version").Parse(greeting))
	if err := tmpl.Execute(w, info); err != nil {
		slog.Error("Failed to print build info", "error", err)
	}
}

// Info возвращает строку с информацией о сборке
func Info() string {
	builder := &strings.Builder{}
	WriteBuildInfo(builder)
	return builder.String()
}
