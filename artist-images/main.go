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

var articles = map[string]bool{
	"the": true, "los": true, "las": true,
	"le": true, "les": true, "el": true, "la": true,
	"die": true, "das": true, "der": true, "gli": true, "il": true,
}

var accentMap = strings.NewReplacer(
	"à", "a", "á", "a", "â", "a", "ã", "a", "ä", "a", "å", "a",
	"è", "e", "é", "e", "ê", "e", "ë", "e",
	"ì", "i", "í", "i", "î", "i", "ï", "i",
	"ò", "o", "ó", "o", "ô", "o", "õ", "o", "ö", "o", "ø", "o",
	"ù", "u", "ú", "u", "û", "u", "ü", "u",
	"ý", "y", "ÿ", "y",
	"ñ", "n", "ç", "c", "ß", "ss",
	"À", "a", "Á", "a", "Â", "a", "Ã", "a", "Ä", "a", "Å", "a",
	"È", "e", "É", "e", "Ê", "e", "Ë", "e",
	"Ì", "i", "Í", "i", "Î", "i", "Ï", "i",
	"Ò", "o", "Ó", "o", "Ô", "o", "Õ", "o", "Ö", "o", "Ø", "o",
	"Ù", "u", "Ú", "u", "Û", "u", "Ü", "u",
	"Ý", "y", "Ñ", "n", "Ç", "c",
)

func normalize(s string) string {
	s = accentMap.Replace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", "and")
	// split on non-alnum so we can detect article boundaries
	parts := nonAlnum.Split(s, -1)
	filtered := parts[:0]
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	// strip leading article only when more words follow
	if len(filtered) > 1 && articles[filtered[0]] {
		filtered = filtered[1:]
	}
	return strings.Join(filtered, "")
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
