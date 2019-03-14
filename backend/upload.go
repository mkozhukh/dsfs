package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

type uploadResponse struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

type uploadHandler struct {
	render *render.Render
}

func (u uploadHandler) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	//make new name
	time := strconv.FormatInt(time.Now().Unix(), 10)
	newName := fmt.Sprintf("%x", md5.Sum([]byte(handler.Filename+time)))

	newPath := filepath.Join(Config.Folder, newName)
	f, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	size, err := io.Copy(f, file)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("File stored, %s %d [%s]", handler.Filename, size, newName)

	path := "/" + Config.Path + "/" + newName + "/" + handler.Filename

	u.render.JSON(w, http.StatusOK, uploadResponse{
		Status: "server",
		Path:   path,
	})
}

func (u uploadHandler) retrieve(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "realname")
	file := chi.URLParam(r, "file")

	w.Header().Set("Content-Type", "applicaiton/binary")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)

	filePath := filepath.Join(Config.Folder, file)
	http.ServeFile(w, r, filePath)
	log.Printf("File requested: %s", name)
}
