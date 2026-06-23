package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"dashboard/database"
	"dashboard/handlers"
)

func findDir(name string) string {
	candidates := []string{
		name,
		filepath.Join("..", name),
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, name),
			filepath.Join(exeDir, "..", name),
		)
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	abs, _ := filepath.Abs(name)
	return abs
}

func main() {
	database.Init("dashboard.db")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.Dashboard)
	mux.HandleFunc("GET /studies", handlers.Studies)
	mux.HandleFunc("GET /aiwork", handlers.AIWork)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /other", handlers.Other)

	// dashboard widgets
	mux.HandleFunc("POST /todos", handlers.CreateTodo)
	mux.HandleFunc("POST /todos/{id}/toggle", handlers.ToggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)
	mux.HandleFunc("POST /reminder", handlers.SaveReminder)
	mux.HandleFunc("POST /stats/{id}", handlers.UpdateStat)

	// studies
	mux.HandleFunc("POST /classes", handlers.CreateClass)
	mux.HandleFunc("PUT /classes/{id}", handlers.UpdateClass)
	mux.HandleFunc("DELETE /classes/{id}", handlers.DeleteClass)
	mux.HandleFunc("POST /quicklinks", handlers.CreateQuickLink)
	mux.HandleFunc("DELETE /quicklinks/{id}", handlers.DeleteQuickLink)

	// aiwork
	mux.HandleFunc("POST /devprojects", handlers.CreateDevProject)
	mux.HandleFunc("PUT /devprojects/{id}", handlers.UpdateDevProject)
	mux.HandleFunc("DELETE /devprojects/{id}", handlers.DeleteDevProject)
	mux.HandleFunc("POST /knowledge", handlers.CreateKnowledge)
	mux.HandleFunc("PUT /knowledge/{id}", handlers.UpdateKnowledge)
	mux.HandleFunc("DELETE /knowledge/{id}", handlers.DeleteKnowledge)
	mux.HandleFunc("POST /curiosities", handlers.CreateCuriosity)
	mux.HandleFunc("PUT /curiosities/{id}", handlers.UpdateCuriosity)
	mux.HandleFunc("DELETE /curiosities/{id}", handlers.DeleteCuriosity)
	mux.HandleFunc("POST /jobapps", handlers.CreateJobApp)
	mux.HandleFunc("PUT /jobapps/{id}", handlers.UpdateJobApp)
	mux.HandleFunc("DELETE /jobapps/{id}", handlers.DeleteJobApp)
	mux.HandleFunc("POST /futureideas", handlers.SaveFutureIdeas)

	// health
	mux.HandleFunc("POST /dietlog", handlers.SaveDietLog)
	mux.HandleFunc("PUT /gymnotes/{id}", handlers.UpdateGymDay)
	mux.HandleFunc("PUT /gymnotes/{id}/toggle", handlers.ToggleGymDay)
	mux.HandleFunc("PUT /skincare/{field}", handlers.ToggleSkincare)
	mux.HandleFunc("POST /hairnotes", handlers.SaveHairNotes)
	mux.HandleFunc("POST /weeklyscores", handlers.CreateWeeklyScore)
	mux.HandleFunc("PUT /weeklyscores/{id}", handlers.UpdateWeeklyScore)
	mux.HandleFunc("DELETE /weeklyscores/{id}", handlers.DeleteWeeklyScore)

	// other
	mux.HandleFunc("POST /othernotes", handlers.SaveOtherNotes)
	mux.HandleFunc("POST /bookmarks", handlers.CreateBookmark)
	mux.HandleFunc("DELETE /bookmarks/{id}", handlers.DeleteBookmark)

	// Serve assets from the correct directory regardless of where we run from
	assetsDir := findDir("assets")
	log.Println("Assets dir:", assetsDir)
	fs := http.FileServer(http.Dir(assetsDir))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
