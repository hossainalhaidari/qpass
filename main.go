package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/bradfitz/gomemcache/memcache"
)

type TemplateData struct {
	Display bool
	Text    string
	Secret  string
}

func main() {
	cacheHost := os.Getenv("MEMCACHED_HOST")
	if cacheHost == "" {
		cacheHost = "localhost:11211"
	}
	cache := memcache.New(cacheHost)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			path := r.URL.Path[1:]
			tmpl, _ := template.ParseFiles("index.html")

			if path == "" {
				tmpl.ExecuteTemplate(w, "layout", nil)
				return
			}

			parts := strings.Split(path, "/")
			if len(parts) != 2 {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "Not Found")
				return
			}

			text := get(cache, parts[0])
			if text == "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "Not Found")
				return
			}

			data := TemplateData{
				Display: true,
				Text:    text,
				Secret:  parts[1],
			}

			tmpl.ExecuteTemplate(w, "layout", data)
		case "POST":
			defer r.Body.Close()
			body, _ := io.ReadAll(r.Body)
			fmt.Fprint(w, set(cache, body))
		}
	})
	http.ListenAndServe(":3000", nil)
}

func set(cache *memcache.Client, body []byte) string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	id := hex.EncodeToString(bytes)

	err := cache.Set(&memcache.Item{Key: id, Value: body, Expiration: 86400})
	if err != nil {
		return ""
	}

	return id
}

func get(cache *memcache.Client, id string) string {
	val, err := cache.Get(id)
	if err != nil {
		return ""
	}

	err = cache.Delete(id)
	if err != nil {
		return ""
	}

	return string(val.Value)
}
