package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tmpl *template.Template

func init() {
	rootDir := "./templates"
	funcMap := template.FuncMap{}
	var err error
	tmpl, err = findAndParseTemplates(rootDir, funcMap)
	if err != nil {
		log.Fatal(err)
	}
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := os.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

// renderer function (handles different layouts)
func renderTemplate(w http.ResponseWriter, layoutName, tmplName string, data interface{}) {
	// Execute the specific template first and get the rendered content
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the layoutData struct to pass to the layout template
	layoutData := struct {
		Content     template.HTML
		Data        interface{}
		UserContent interface{}
	}{
		Content:     template.HTML(buf.String()),
		Data:        data,
		UserContent: data,
	}

	// Execute the layout template, passing the layoutData struct as the template data
	err = tmpl.ExecuteTemplate(w, layoutName, layoutData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
