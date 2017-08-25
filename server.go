package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"mkozhukh/dsfs/auth"
	"mkozhukh/dsfs/config"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type contextKey int

type uploadResponse struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

const (
	accessType contextKey = iota
)

var format = render.New(render.Options{
	Asset:      Asset,
	AssetNames: AssetNames,
	Directory:  "static",
	Extensions: []string{".html"},
	//IsDevelopment: true,
})

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	path := "config.yml"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	config.Config.Load(path)

	// Routes
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(sessionMiddleware())
	r.Use(auth.ForceLogin)

	r.Mount("/debug", middleware.Profiler())
	//add login|logout routes
	r.Mount("/auth", auth.AddRoutes(format))

	//get list of files
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("upload")
		if err != nil {
			log.Error(err)
			return
		}
		defer file.Close()

		//make new name
		time := strconv.FormatInt(time.Now().Unix(), 10)
		newName := fmt.Sprintf("%x", md5.Sum([]byte(handler.Filename+time)))

		newPath := filepath.Join(config.Config.Folder, newName)
		f, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Error(err)
			return
		}
		defer f.Close()

		size, err := io.Copy(f, file)
		if err != nil {
			log.Error(err)
			return
		}

		log.WithFields(log.Fields{
			"name": handler.Filename,
			"size": size,
			"id":   newName,
		}).Info("File stored")

		path := "/" + config.Config.Path + "/" + newName + "/" + handler.Filename
		format.JSON(w, http.StatusOK, uploadResponse{
			Status: "server",
			Path:   path,
		})
	})

	//serve uploaded files
	r.Get("/"+config.Config.Path+"/{file}/{realname}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "realname")
		file := chi.URLParam(r, "file")

		w.Header().Set("Content-Type", "applicaiton/binary")
		w.Header().Set("Content-Disposition", "attachment; filename="+name)

		filePath := filepath.Join(config.Config.Folder, file)
		http.ServeFile(w, r, filePath)
		log.WithField("name", name).Info("File request")
	})

	//static assets
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		format.HTML(w, http.StatusOK, "index", nil)
	})

	http.ListenAndServe(":"+config.Config.Server.Port, r)
}
