package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"dashboard/database"
)

type Class struct {
	ID            int64
	Name          string
	MaterialEmoji string
	ExamDates     string
	HwDeadline    string
}

type QuickLink struct {
	ID    int64
	Label string
	URL   string
}

type StudiesData struct {
	Classes    []Class
	QuickLinks []QuickLink
}

var materialEmojis = []string{"📘", "🧪", "🧮", "🎨", "🌍", "💻", "📐", "🧬"}

func Studies(w http.ResponseWriter, r *http.Request) {
	data := StudiesData{
		Classes:    fetchClasses(),
		QuickLinks: fetchQuickLinks("studies"),
	}
	renderPage(w, r, "studies", data)
}

func fetchClasses() []Class {
	rows, err := database.DB.Query("SELECT id, name, material_emoji, exam_dates, hw_deadline FROM classes ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var classes []Class
	for rows.Next() {
		var c Class
		rows.Scan(&c.ID, &c.Name, &c.MaterialEmoji, &c.ExamDates, &c.HwDeadline)
		classes = append(classes, c)
	}
	return classes
}

func fetchQuickLinks(page string) []QuickLink {
	rows, err := database.DB.Query("SELECT id, label, url FROM quick_links WHERE page = ? ORDER BY id", page)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var links []QuickLink
	for rows.Next() {
		var l QuickLink
		rows.Scan(&l.ID, &l.Label, &l.URL)
		links = append(links, l)
	}
	return links
}

type classRowData struct {
	Class
	Emojis []string
}

var classRowTmpl = template.Must(template.New("classRow").Parse(`
<tr id="class-{{.ID}}">
	<td><input type="text" value="{{.Name}}" name="name" hx-put="/classes/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Class name"></td>
	<td>
		<select name="material_emoji" hx-put="/classes/{{.ID}}" hx-include="closest tr" hx-trigger="change" hx-swap="none">
			{{$cur := .MaterialEmoji}}
			{{range .Emojis}}<option value="{{.}}" {{if eq . $cur}}selected{{end}}>{{.}}</option>{{end}}
		</select>
	</td>
	<td><input type="text" value="{{.ExamDates}}" name="exam_dates" hx-put="/classes/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Exam dates"></td>
	<td><input type="text" value="{{.HwDeadline}}" name="hw_deadline" hx-put="/classes/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="HW deadline"></td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/classes/{{.ID}}" hx-target="#class-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func renderClassRow(w http.ResponseWriter, c Class) {
	renderFragment(w, classRowTmpl, classRowData{Class: c, Emojis: materialEmojis})
}

func CreateClass(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO classes (name, material_emoji, exam_dates, hw_deadline) VALUES ('', '📘', '', '')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderClassRow(w, Class{ID: id, MaterialEmoji: "📘"})
}

func UpdateClass(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE classes SET name = ?, material_emoji = ?, exam_dates = ?, hw_deadline = ? WHERE id = ?",
		r.FormValue("name"), r.FormValue("material_emoji"), r.FormValue("exam_dates"), r.FormValue("hw_deadline"), id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteClass(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM classes WHERE id = ?", id)
	w.Write([]byte(""))
}

var quickLinkTmpl = template.Must(template.New("quickLink").Parse(`
<div class="quick-link-chip" id="quicklink-{{.ID}}">
	<a href="{{.URL}}" target="_blank" rel="noopener">{{.Label}}</a>
	<span class="x" hx-delete="/quicklinks/{{.ID}}" hx-target="#quicklink-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</span>
</div>`))

func CreateQuickLink(w http.ResponseWriter, r *http.Request) {
	label := r.FormValue("label")
	url := r.FormValue("url")
	page := r.FormValue("page")
	if page == "" {
		page = "studies"
	}
	if label == "" || url == "" {
		w.Write([]byte(""))
		return
	}
	res, err := database.DB.Exec("INSERT INTO quick_links (label, url, page) VALUES (?, ?, ?)", label, url, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, quickLinkTmpl, QuickLink{ID: id, Label: label, URL: url})
}

func DeleteQuickLink(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM quick_links WHERE id = ?", id)
	w.Write([]byte(""))
}
