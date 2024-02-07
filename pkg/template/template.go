package template

import (
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// template helper from file and cache
type Template struct {
	prefixPath string
	templates  map[string]*template.Template
}

// NewTemplate create new template instance
func New(prefixPath string) *Template {

	cwd, err := os.Getwd()
	if err != nil {
		log.Error("Error getting current working directory", "error", err)
		panic(err)
	}

	rootPath := filepath.Clean(cwd + "/../..")

	return &Template{
		prefixPath: fmt.Sprintf("%s/%s", rootPath, prefixPath),
		templates:  make(map[string]*template.Template),
	}
}

func (t *Template) Render(filename string, w http.ResponseWriter, data interface{}) {

	if _, ok := t.templates[filename]; !ok {

		fullPath := fmt.Sprintf("%s/%s", t.prefixPath, filename)
		log.Info("template", "path", fullPath)
		tpl, err := template.ParseFiles(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.templates[filename] = tpl

	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	t.templates[filename].Execute(w, data)

}
