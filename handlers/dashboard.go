package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"dashboard/database"
)

type Todo struct {
	ID   int64
	Text string
	Done bool
}

type Stat struct {
	ID    int64
	Label string
	Value int
}

type DashboardData struct {
	Todos    []Todo
	Reminder string
	Stats    []Stat
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	data := DashboardData{
		Todos:    fetchTodos(),
		Reminder: fetchReminder(),
		Stats:    fetchStats(),
	}
	renderPage(w, r, "dashboard", data)
}

func fetchTodos() []Todo {
	rows, err := database.DB.Query("SELECT id, text, done FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var todos []Todo
	for rows.Next() {
		var t Todo
		rows.Scan(&t.ID, &t.Text, &t.Done)
		todos = append(todos, t)
	}
	return todos
}

func fetchReminder() string {
	var content string
	database.DB.QueryRow("SELECT content FROM reminder_text ORDER BY id LIMIT 1").Scan(&content)
	return content
}

func fetchStats() []Stat {
	rows, err := database.DB.Query("SELECT id, label, value FROM stats ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var stats []Stat
	for rows.Next() {
		var s Stat
		rows.Scan(&s.ID, &s.Label, &s.Value)
		stats = append(stats, s)
	}
	return stats
}

var todoItemTmpl = template.Must(template.New("todoItem").Parse(`
<li class="{{if .Done}}done{{end}}" id="todo-{{.ID}}">
	<input type="checkbox" class="todo-check" {{if .Done}}checked{{end}}
		hx-post="/todos/{{.ID}}/toggle" hx-target="#todo-{{.ID}}" hx-swap="outerHTML">
	<span class="todo-text">{{.Text}}</span>
	<button type="button" class="btn-small btn-danger" hx-delete="/todos/{{.ID}}" hx-target="#todo-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button>
</li>`))

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		w.Write([]byte(""))
		return
	}
	res, err := database.DB.Exec("INSERT INTO todos (text, done) VALUES (?, 0)", text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, todoItemTmpl, Todo{ID: id, Text: text, Done: false})
}

func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE todos SET done = NOT done WHERE id = ?", id)
	var t Todo
	t.ID = id
	database.DB.QueryRow("SELECT text, done FROM todos WHERE id = ?", id).Scan(&t.Text, &t.Done)
	renderFragment(w, todoItemTmpl, t)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	w.Write([]byte(""))
}

var reminderTmpl = template.Must(template.New("reminder").Parse(`
<input type="text" id="reminder-input" name="content" value="{{.}}"
	hx-post="/reminder" hx-target="#reminder-input" hx-swap="outerHTML" hx-trigger="blur">`))

func SaveReminder(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	database.DB.Exec("UPDATE reminder_text SET content = ? WHERE id = (SELECT id FROM reminder_text ORDER BY id LIMIT 1)", content)
	renderFragment(w, reminderTmpl, content)
}

func UpdateStat(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	value, _ := strconv.Atoi(r.FormValue("value"))
	database.DB.Exec("UPDATE stats SET value = ? WHERE id = ?", value, id)
	w.WriteHeader(http.StatusNoContent)
}
