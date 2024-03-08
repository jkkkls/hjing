package web_backend

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/utils"
)

//go:embed dist
var reactStatic embed.FS
var uiFS fs.FS

func init() {
	var err error
	uiFS, err = fs.Sub(reactStatic, "dist")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}
}

func initEmbedReact() {
	addr := config.GetString("web", "webAddress")
	if addr == "" {
		return
	}

	h := http.NewServeMux()

	h.HandleFunc("/api/", handleApi)
	h.HandleFunc("/", handleStatic)

	utils.Go(func() {
		err := http.ListenAndServe(addr, h)
		if err != nil {
			panic(err)
		}
	})
}

func handleApi(w http.ResponseWriter, r *http.Request) {
	// 代理到后端9094端口
	u, _ := url.Parse(fmt.Sprintf("http://127.0.0.1%v/", config.GetString("web", "address")))
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.ServeHTTP(w, r)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := filepath.Clean(r.URL.Path)
	if path == "/" { // Add other paths that you route on the UI side here
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	file, err := uiFS.Open(path)
	if err != nil {
		// try_files $uri $uri/ /index.html;
		path = "index.html"
		file, err = uiFS.Open(path)
	}
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(w, r)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}
