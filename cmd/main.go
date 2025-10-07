package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/handlers"
)

func main() {
	timerHandler := handlers.TimerHandler{}
	// fileServer := http.FileServer(http.Dir("./assets/"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", assetsHandler)
	mux.HandleFunc("/createTimer", timerHandler.CreateTimer)
	mux.HandleFunc("/pauseTimer", timerHandler.PauseTimer)
	mux.HandleFunc("/startTimer", timerHandler.StartTimer)
	mux.HandleFunc("/deleteTimer", timerHandler.DeleteTimer)
	mux.HandleFunc("/createUser", timerHandler.CreateUser)
	mux.HandleFunc("/updateTimerTittle", timerHandler.UpdateTimerTitle)
	mux.HandleFunc("/", timerHandler.RenderUserTimers)

	log.Fatal(http.ListenAndServe(":"+global.PORT, mux))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.Local = time.UTC
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/assets/")

	// Устанавливаем Content-Type в зависимости от расширения
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	// Кэширование
	if ext == ".css" || ext == ".js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}

	// Отдаем файл
	http.ServeFile(w, r, "./assets/"+path)
}
