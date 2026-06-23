package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"dashboard/database"
)

type Bookmark struct {
	ID    int64
	Title string
	URL   string
}

type OtherData struct {
	Notes     string
	Bookmarks []Bookmark
}

func Other(w http.ResponseWriter, r *http.Request) {
	data := OtherData{
		Notes:     fetchOtherNotes(),
		Bookmarks: fetchBookmarks(),
	}
	renderPage(w, r, "other", data)
}

func fetchOtherNotes() string {
	var content string
	database.DB.QueryRow("SELECT content FROM other_notes ORDER BY id LIMIT 1").Scan(&content)
	return content
}

func SaveOtherNotes(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	database.DB.Exec("UPDATE other_notes SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = (SELECT id FROM other_notes ORDER BY id LIMIT 1)", content)
	w.WriteHeader(http.StatusNoContent)
}

func fetchBookmarks() []Bookmark {
	rows, err := database.DB.Query("SELECT id, title, url FROM bookmarks ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []Bookmark
	for rows.Next() {
		var b Bookmark
		rows.Scan(&b.ID, &b.Title, &b.URL)
		list = append(list, b)
	}
	return list
}

var bookmarkItemTmpl = template.Must(template.New("bookmarkItem").Parse(`
<li id="bookmark-{{.ID}}">
	<a href="{{.URL}}" target="_blank" rel="noopener">{{.Title}}</a>
	<button type="button" class="btn-small btn-danger" hx-delete="/bookmarks/{{.ID}}" hx-target="#bookmark-{{.ID}}" hx-swap="outerHTML swap:0.15s">✕</button>
</li>`))

func CreateBookmark(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	url := r.FormValue("url")
	if title == "" || url == "" {
		w.Write([]byte(""))
		return
	}
	res, err := database.DB.Exec("INSERT INTO bookmarks (title, url) VALUES (?, ?)", title, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	renderFragment(w, bookmarkItemTmpl, Bookmark{ID: id, Title: title, URL: url})
}

func DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	database.DB.Exec("DELETE FROM bookmarks WHERE id = ?", id)
	w.Write([]byte(""))
}
