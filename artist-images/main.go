package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var coverMap = map[string]string{}
var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
var logRequests bool

func normalize(s string) string {
	s = strings.ToLower(s)
	s = nonAlnum.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	return s
}

func indexDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("cannot read %q: %v", dir, err)
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		cover := filepath.Join(dir, e.Name(), "cover.jpg")
		if _, err := os.Stat(cover); err == nil {
			// Artist folder — has cover.jpg directly inside
			coverMap[normalize(e.Name())] = cover
		} else {
			// Partition folder (e.g. "1/") — scan one level deeper
			indexDir(filepath.Join(dir, e.Name()))
		}
	}
}

func buildMap(root string) {
	indexDir(root)
	log.Printf("indexed %d artist covers from %s", len(coverMap), root)
}

func main() {
	root := os.Getenv("MUSIC_ROOT")
	if root == "" {
		root = "/media/music"
	}
	logRequests = strings.ToLower(os.Getenv("LOG_REQUESTS")) == "true"
	buildMap(root)

	http.HandleFunc("/avatar", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing name param", http.StatusBadRequest)
			return
		}
		path, ok := coverMap[normalize(name)]
		if !ok {
			if logRequests {
				log.Printf("MISS  %q (%.0fms)", name, float64(time.Since(start).Microseconds())/1000)
			}
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, path)
		if logRequests {
			log.Printf("HIT   %q (%.0fms)", name, float64(time.Since(start).Microseconds())/1000)
		}
	})

	log.Printf("listening on :8080 (LOG_REQUESTS=%v)", logRequests)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
