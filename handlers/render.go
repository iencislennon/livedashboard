package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var pages = map[string]*template.Template{}

func init() {
	// Find templates relative to the executable or working directory
	baseDir := findTemplatesDir()
	log.Println("Templates dir:", baseDir)

	pageFiles := map[string]string{
		"dashboard": "dashboard.html",
		"studies":   "studies.html",
		"aiwork":    "aiwork.html",
		"health":    "health.html",
		"other":     "other.html",
	}
	baseTmpl := filepath.Join(baseDir, "base.html")
	for name, file := range pageFiles {
		path := filepath.Join(baseDir, file)
		t, err := template.ParseFiles(baseTmpl, path)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
		pages[name] = t
	}
}

func findTemplatesDir() string {
	// Try relative to working directory first
	candidates := []string{
		"templates",
		"../templates",
	}
	// Also try relative to executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "templates"),
			filepath.Join(exeDir, "../templates"),
		)
	}
	for _, c := range candidates {
		if info, err := os.Stat(filepath.Join(c, "base.html")); err == nil && !info.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	// Default fallback
	abs, _ := filepath.Abs("templates")
	return abs
}

// renderPage renders the full document on a direct GET, or just the
// "content" block when the request came from HTMX navigation swap.
func renderPage(w http.ResponseWriter, r *http.Request, name string, data any) {
	tmpl, ok := pages[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var err error
	if r.Header.Get("HX-Request") == "true" {
		err = tmpl.ExecuteTemplate(w, "content", data)
	} else {
		err = tmpl.ExecuteTemplate(w, "base", data)
	}
	if err != nil {
		log.Printf("Template render error (%s): %v", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderFragment(w http.ResponseWriter, tmpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Fragment render error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
