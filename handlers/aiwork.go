package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"dashboard/database"
)

type DevProject struct {
	ID        int64
	Name      string
	RepoURL   string
	ClaudeURL string
	OtherURL  string
}

type Knowledge struct {
	ID    int64
	Title string
	Topic string
	Link  string
}

type Curiosity struct {
	ID          int64
	Name        string
	Link        string
	Description string
	Notes       string
	IconEmoji   string
}

type JobApp struct {
	ID       int64
	Company  string
	Title    string
	AppURL   string
	CVURL    string
	Priority string
}

type AIWorkData struct {
	DevProjects  []DevProject
	Knowledge    []Knowledge
	Curiosities  []Curiosity
	JobApps      []JobApp
	FutureIdeas  string
}

var curiosityIcons = []string{"✨", "🔮", "🧙", "🗝️", "🛰️", "🔬", "🎲", "📜"}

func AIWork(w http.ResponseWriter, r *http.Request) {
	data := AIWorkData{
		DevProjects: fetchDevProjects(),
		Knowledge:   fetchKnowledge(),
		Curiosities: fetchCuriosities(),
		JobApps:     fetchJobApps(),
		FutureIdeas: fetchFutureIdeas(),
	}
	renderPage(w, r, "aiwork", data)
}

func fetchDevProjects() []DevProject {
	rows, err := database.DB.Query("SELECT id, name, repo_url, claude_url, other_url FROM dev_projects ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []DevProject
	for rows.Next() {
		var d DevProject
		rows.Scan(&d.ID, &d.Name, &d.RepoURL, &d.ClaudeURL, &d.OtherURL)
		list = append(list, d)
	}
	return list
}

func fetchKnowledge() []Knowledge {
	rows, err := database.DB.Query("SELECT id, title, topic, link FROM knowledge ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []Knowledge
	for rows.Next() {
		var k Knowledge
		rows.Scan(&k.ID, &k.Title, &k.Topic, &k.Link)
		list = append(list, k)
	}
	return list
}

func fetchCuriosities() []Curiosity {
	rows, err := database.DB.Query("SELECT id, name, link, description, notes, icon_emoji FROM curiosities ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []Curiosity
	for rows.Next() {
		var c Curiosity
		rows.Scan(&c.ID, &c.Name, &c.Link, &c.Description, &c.Notes, &c.IconEmoji)
		list = append(list, c)
	}
	return list
}

func fetchJobApps() []JobApp {
	rows, err := database.DB.Query("SELECT id, company, title, app_url, cv_url, priority FROM job_apps ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []JobApp
	for rows.Next() {
		var j JobApp
		rows.Scan(&j.ID, &j.Company, &j.Title, &j.AppURL, &j.CVURL, &j.Priority)
		list = append(list, j)
	}
	return list
}

func fetchFutureIdeas() string {
	var content string
	database.DB.QueryRow("SELECT content FROM future_ideas ORDER BY id LIMIT 1").Scan(&content)
	return content
}

// ---------- dev_projects ----------

var devProjectRowTmpl = template.Must(template.New("devProjectRow").Parse(`
<tr id="devproject-{{.ID}}">
	<td><input type="text" value="{{.Name}}" name="name" hx-put="/devprojects/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Name"></td>
	<td><input type="url" value="{{.RepoURL}}" name="repo_url" hx-put="/devprojects/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Repo link"></td>
	<td><input type="url" value="{{.ClaudeURL}}" name="claude_url" hx-put="/devprojects/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Claude chat link"></td>
	<td><input type="url" value="{{.OtherURL}}" name="other_url" hx-put="/devprojects/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="URL"></td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/devprojects/{{.ID}}" hx-target="#devproject-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func CreateDevProject(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO dev_projects (name, repo_url, claude_url, other_url) VALUES ('', '', '', '')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, devProjectRowTmpl, DevProject{ID: id})
}

func UpdateDevProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE dev_projects SET name = ?, repo_url = ?, claude_url = ?, other_url = ? WHERE id = ?",
		r.FormValue("name"), r.FormValue("repo_url"), r.FormValue("claude_url"), r.FormValue("other_url"), id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteDevProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM dev_projects WHERE id = ?", id)
	w.Write([]byte(""))
}

// ---------- knowledge ----------

var knowledgeRowTmpl = template.Must(template.New("knowledgeRow").Parse(`
<tr id="knowledge-{{.ID}}">
	<td><input type="text" value="{{.Title}}" name="title" hx-put="/knowledge/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Title"></td>
	<td><input type="text" value="{{.Topic}}" name="topic" hx-put="/knowledge/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Topic"></td>
	<td><input type="url" value="{{.Link}}" name="link" hx-put="/knowledge/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Link"></td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/knowledge/{{.ID}}" hx-target="#knowledge-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func CreateKnowledge(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO knowledge (title, topic, link) VALUES ('', '', '')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, knowledgeRowTmpl, Knowledge{ID: id})
}

func UpdateKnowledge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE knowledge SET title = ?, topic = ?, link = ? WHERE id = ?",
		r.FormValue("title"), r.FormValue("topic"), r.FormValue("link"), id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteKnowledge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM knowledge WHERE id = ?", id)
	w.Write([]byte(""))
}

// ---------- curiosities ----------

type curiosityRowData struct {
	Curiosity
	Icons []string
}

var curiosityRowTmpl = template.Must(template.New("curiosityRow").Parse(`
<tr id="curiosity-{{.ID}}">
	<td>
		<select name="icon_emoji" hx-put="/curiosities/{{.ID}}" hx-include="closest tr" hx-trigger="change" hx-swap="none">
			{{$cur := .IconEmoji}}
			{{range .Icons}}<option value="{{.}}" {{if eq . $cur}}selected{{end}}>{{.}}</option>{{end}}
		</select>
	</td>
	<td><input type="text" value="{{.Name}}" name="name" hx-put="/curiosities/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Name"></td>
	<td><input type="url" value="{{.Link}}" name="link" hx-put="/curiosities/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Link"></td>
	<td><input type="text" value="{{.Description}}" name="description" hx-put="/curiosities/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Description"></td>
	<td><input type="text" value="{{.Notes}}" name="notes" hx-put="/curiosities/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Notes"></td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/curiosities/{{.ID}}" hx-target="#curiosity-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func renderCuriosityRow(w http.ResponseWriter, c Curiosity) {
	renderFragment(w, curiosityRowTmpl, curiosityRowData{Curiosity: c, Icons: curiosityIcons})
}

func CreateCuriosity(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO curiosities (name, link, description, notes, icon_emoji) VALUES ('', '', '', '', '✨')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderCuriosityRow(w, Curiosity{ID: id, IconEmoji: "✨"})
}

func UpdateCuriosity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE curiosities SET name = ?, link = ?, description = ?, notes = ?, icon_emoji = ? WHERE id = ?",
		r.FormValue("name"), r.FormValue("link"), r.FormValue("description"), r.FormValue("notes"), r.FormValue("icon_emoji"), id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteCuriosity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM curiosities WHERE id = ?", id)
	w.Write([]byte(""))
}

// ---------- job_apps ----------

var jobAppRowTmpl = template.Must(template.New("jobAppRow").Parse(`
<tr id="jobapp-{{.ID}}">
	<td><input type="text" value="{{.Company}}" name="company" hx-put="/jobapps/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Company"></td>
	<td><input type="text" value="{{.Title}}" name="title" hx-put="/jobapps/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Title"></td>
	<td><input type="url" value="{{.AppURL}}" name="app_url" hx-put="/jobapps/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Application link"></td>
	<td><input type="url" value="{{.CVURL}}" name="cv_url" hx-put="/jobapps/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="CV link"></td>
	<td>
		<select name="priority" hx-put="/jobapps/{{.ID}}" hx-include="closest tr" hx-trigger="change" hx-target="closest tr" hx-swap="outerHTML">
			<option value="high" {{if eq .Priority "high"}}selected{{end}}>High</option>
			<option value="medium" {{if eq .Priority "medium"}}selected{{end}}>Medium</option>
			<option value="low" {{if eq .Priority "low"}}selected{{end}}>Low</option>
		</select>
		<span class="gem gem-{{.Priority}}"></span>
	</td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/jobapps/{{.ID}}" hx-target="#jobapp-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func CreateJobApp(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO job_apps (company, title, app_url, cv_url, priority) VALUES ('', '', '', '', 'medium')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, jobAppRowTmpl, JobApp{ID: id, Priority: "medium"})
}

func UpdateJobApp(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	priority := r.FormValue("priority")
	database.DB.Exec("UPDATE job_apps SET company = ?, title = ?, app_url = ?, cv_url = ?, priority = ? WHERE id = ?",
		r.FormValue("company"), r.FormValue("title"), r.FormValue("app_url"), r.FormValue("cv_url"), priority, id)

	var j JobApp
	j.ID = id
	database.DB.QueryRow("SELECT company, title, app_url, cv_url, priority FROM job_apps WHERE id = ?", id).
		Scan(&j.Company, &j.Title, &j.AppURL, &j.CVURL, &j.Priority)
	renderFragment(w, jobAppRowTmpl, j)
}

func DeleteJobApp(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM job_apps WHERE id = ?", id)
	w.Write([]byte(""))
}

// ---------- future ideas ----------

func SaveFutureIdeas(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	database.DB.Exec("UPDATE future_ideas SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = (SELECT id FROM future_ideas ORDER BY id LIMIT 1)", content)
	w.WriteHeader(http.StatusNoContent)
}
