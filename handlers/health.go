package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"dashboard/database"
)

type GymNote struct {
	ID        int64
	DayName   string
	Exercises string
	Active    bool
}

type WeeklyScore struct {
	ID        int64
	WeekLabel string
	Score     int
}

type HealthData struct {
	Today        string
	DietContent  string
	GymNotes     []GymNote
	Skincare     SkincareState
	HairNotes    string
	WeeklyScores []WeeklyScore
}

type SkincareState struct {
	Cleanser    bool
	Moisturizer bool
	SPF         bool
	Serum       bool
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func Health(w http.ResponseWriter, r *http.Request) {
	t := today()
	data := HealthData{
		Today:        t,
		DietContent:  fetchDietLog(t),
		GymNotes:     fetchGymNotes(),
		Skincare:     fetchSkincare(t),
		HairNotes:    fetchHairNotes(),
		WeeklyScores: fetchWeeklyScores(),
	}
	renderPage(w, r, "health", data)
}

func fetchDietLog(date string) string {
	var content string
	database.DB.QueryRow("SELECT content FROM diet_log WHERE date = ?", date).Scan(&content)
	return content
}

func SaveDietLog(w http.ResponseWriter, r *http.Request) {
	t := today()
	content := r.FormValue("content")
	var id int64
	err := database.DB.QueryRow("SELECT id FROM diet_log WHERE date = ?", t).Scan(&id)
	if err != nil {
		database.DB.Exec("INSERT INTO diet_log (date, content) VALUES (?, ?)", t, content)
	} else {
		database.DB.Exec("UPDATE diet_log SET content = ? WHERE id = ?", content, id)
	}
	w.WriteHeader(http.StatusNoContent)
}

func fetchGymNotes() []GymNote {
	rows, err := database.DB.Query("SELECT id, day_name, exercises, active FROM gym_notes ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []GymNote
	for rows.Next() {
		var g GymNote
		rows.Scan(&g.ID, &g.DayName, &g.Exercises, &g.Active)
		list = append(list, g)
	}
	return list
}

var gymDayTmpl = template.Must(template.New("gymDay").Parse(`
<div class="gym-day" id="gym-{{.ID}}">
	<div class="gym-day-header">
		<span class="day-name">{{.DayName}}</span>
		<div class="toggle-switch {{if .Active}}on{{end}}" hx-put="/gymnotes/{{.ID}}/toggle" hx-target="#gym-{{.ID}}" hx-swap="outerHTML">
			<div class="knob"></div>
		</div>
	</div>
	{{if .Active}}
	<div class="gym-exercises">
		<textarea rows="3" name="exercises" placeholder="Exercises, sets x reps..."
			hx-put="/gymnotes/{{.ID}}" hx-trigger="blur" hx-swap="none">{{.Exercises}}</textarea>
	</div>
	{{end}}
</div>`))

func ToggleGymDay(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE gym_notes SET active = NOT active WHERE id = ?", id)
	var g GymNote
	g.ID = id
	database.DB.QueryRow("SELECT day_name, exercises, active FROM gym_notes WHERE id = ?", id).Scan(&g.DayName, &g.Exercises, &g.Active)
	renderFragment(w, gymDayTmpl, g)
}

func UpdateGymDay(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("UPDATE gym_notes SET exercises = ? WHERE id = ?", r.FormValue("exercises"), id)
	w.WriteHeader(http.StatusNoContent)
}

func fetchSkincare(date string) SkincareState {
	var s SkincareState
	database.DB.QueryRow("SELECT cleanser, moisturizer, spf, serum FROM skincare_log WHERE date = ?", date).
		Scan(&s.Cleanser, &s.Moisturizer, &s.SPF, &s.Serum)
	return s
}

func ToggleSkincare(w http.ResponseWriter, r *http.Request) {
	field := r.PathValue("field")
	allowed := map[string]bool{"cleanser": true, "moisturizer": true, "spf": true, "serum": true}
	if !allowed[field] {
		http.Error(w, "invalid field", http.StatusBadRequest)
		return
	}
	t := today()
	var id int64
	err := database.DB.QueryRow("SELECT id FROM skincare_log WHERE date = ?", t).Scan(&id)
	if err != nil {
		database.DB.Exec("INSERT INTO skincare_log (date, cleanser, moisturizer, spf, serum) VALUES (?, 0, 0, 0, 0)", t)
		database.DB.QueryRow("SELECT id FROM skincare_log WHERE date = ?", t).Scan(&id)
	}
	database.DB.Exec("UPDATE skincare_log SET "+field+" = NOT "+field+" WHERE id = ?", id)
	w.WriteHeader(http.StatusNoContent)
}

func fetchHairNotes() string {
	var content string
	database.DB.QueryRow("SELECT content FROM hair_notes ORDER BY id LIMIT 1").Scan(&content)
	return content
}

func SaveHairNotes(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	database.DB.Exec("UPDATE hair_notes SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = (SELECT id FROM hair_notes ORDER BY id LIMIT 1)", content)
	w.WriteHeader(http.StatusNoContent)
}

func fetchWeeklyScores() []WeeklyScore {
	rows, err := database.DB.Query("SELECT id, week_label, score FROM weekly_scores ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []WeeklyScore
	for rows.Next() {
		var s WeeklyScore
		rows.Scan(&s.ID, &s.WeekLabel, &s.Score)
		list = append(list, s)
	}
	return list
}

var weeklyScoreRowTmpl = template.Must(template.New("weeklyScoreRow").Parse(`
<tr id="weeklyscore-{{.ID}}">
	<td><input type="text" value="{{.WeekLabel}}" name="week_label" hx-put="/weeklyscores/{{.ID}}" hx-include="closest tr" hx-trigger="blur" hx-swap="none" placeholder="Week label"></td>
	<td><input type="number" min="1" max="10" value="{{.Score}}" name="score" hx-put="/weeklyscores/{{.ID}}" hx-include="closest tr" hx-trigger="change" hx-swap="none" onchange="window.dispatchEvent(new Event('weeklyscores-changed'))"></td>
	<td><button type="button" class="btn-small btn-danger" hx-delete="/weeklyscores/{{.ID}}" hx-target="#weeklyscore-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button></td>
</tr>`))

func CreateWeeklyScore(w http.ResponseWriter, r *http.Request) {
	res, err := database.DB.Exec("INSERT INTO weekly_scores (week_label, score) VALUES ('New Week', 5)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, weeklyScoreRowTmpl, WeeklyScore{ID: id, WeekLabel: "New Week", Score: 5})
}

func UpdateWeeklyScore(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	score, _ := strconv.Atoi(r.FormValue("score"))
	database.DB.Exec("UPDATE weekly_scores SET week_label = ?, score = ? WHERE id = ?", r.FormValue("week_label"), score, id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteWeeklyScore(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM weekly_scores WHERE id = ?", id)
	w.Write([]byte(""))
}
