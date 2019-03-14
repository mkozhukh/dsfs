// go:generate go-bindata static/
//go:generate jstore User store/user.go

package main

import (
	"log"
	"net/http"

	// auth
	"github.com/unrolled/render"

	//http server
	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkozhukh/go-remote"
)

func main() {
	Config.Load("")

	var session = scs.NewCookieManager(Config.Server.Secret)
	var output = render.New(render.Options{
		Asset:      Asset,
		AssetNames: AssetNames,
		Directory:  "static",
		Extensions: []string{".html"},
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(auth.GuardRequest("/denied", CanUploadFiles))

	api := remote.NewServer()
	initAPI(api)

	initLogin(r, session)
	initUser(r, api)

	// Routes
	files := uploadHandler{}
	r.Post("/upload", files.upload)                             //store file
	r.Get("/"+Config.Path+"/{file}/{realname}", files.retrieve) //serve uploaded files

	//static assets
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		output.HTML(w, http.StatusOK, "index", nil)
	})

	log.Printf("starting on port %s", Config.Server.Port)
	http.ListenAndServe(":"+Config.Server.Port, r)
}
