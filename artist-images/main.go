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


func buildMap(root string) {
	fresh := map[string]string{}
	letters, err := os.ReadDir(root)
	if err != nil {
		log.Printf("cannot read root %q: %v", root, err)
		coverMap = fresh
		return
	}
	for _, letter := range letters {
		if !letter.IsDir() {
			continue
		}
		artists, err := os.ReadDir(filepath.Join(root, letter.Name()))
		if err != nil {
			log.Printf("cannot read %q: %v", letter.Name(), err)
			continue
		}
		for _, artist := range artists {
			if !artist.IsDir() {
				continue
			}
			cover := filepath.Join(root, letter.Name(), artist.Name(), "cover.jpg")
			if _, err := os.Stat(cover); err == nil {
				fresh[normalize(artist.Name())] = cover
			}
		}
	}
	coverMap = fresh
	log.Printf("indexed %d artist covers from %s", len(coverMap), root)
}

func main() {
	root := os.Getenv("MUSIC_ROOT")
	if root == "" {
		root = "/media/music"
	}
	logRequests = strings.ToLower(os.Getenv("LOG_REQUESTS")) == "true"
	buildMap(root)
	if len(coverMap) == 0 {
		log.Fatalf("no artist covers found in %s — NFS not ready?", root)
	}

	go func() {
		for range time.Tick(time.Hour) {
			log.Printf("rescanning %s", root)
			buildMap(root)
		}
	}()

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
