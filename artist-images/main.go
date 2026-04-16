package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var coverMap = map[string]string{}
var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func normalize(s string) string {
	s = strings.ToLower(s)
	s = nonAlnum.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	return s
}

func buildMap(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("cannot read music root %q: %v", root, err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		cover := filepath.Join(root, e.Name(), "cover.jpg")
		if _, err := os.Stat(cover); err == nil {
			coverMap[normalize(e.Name())] = cover
		}
	}
	log.Printf("indexed %d artist covers from %s", len(coverMap), root)
}

func main() {
	root := os.Getenv("MUSIC_ROOT")
	if root == "" {
		root = "/media/music"
	}
	buildMap(root)

	http.HandleFunc("/avatar", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing name param", http.StatusBadRequest)
			return
		}
		path, ok := coverMap[normalize(name)]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, path)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
