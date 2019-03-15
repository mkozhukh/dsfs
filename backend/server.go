// go:generate go-bindata static/

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
	"github.com/go-chi/cors"

	//modules
	"github.com/mkozhukh/dsfs/backend/api"
	"github.com/mkozhukh/dsfs/backend/auth"
)

func main() {
	Config.Load("")

	var session = scs.NewCookieManager(Config.Server.Secret)
	var output = render.New(render.Options{
		// Asset:      Asset,
		// AssetNames: AssetNames,
		IsDevelopment: true,
		Directory:     "static",
		Extensions:    []string{".html"},
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//cors
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type", "Remote-Csrf"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	// authentication
	initLogin(r, session)
	auth.ReloadRoles()

	//API
	r.Handle("/api/v1", api.New())

	// Routes
	files := uploadHandler{output}
	r.Post("/upload", files.upload)                             //store file
	r.Get("/"+Config.Path+"/{file}/{realname}", files.retrieve) //serve uploaded files

	//static assets
	r.Get("/denied", func(w http.ResponseWriter, r *http.Request) {
		output.HTML(w, http.StatusOK, "denied", nil)
	})

	//client
	rp := r.With(auth.Registry.GuardRequest("/denied", auth.UploadFiles))
	rp.Get("/", func(w http.ResponseWriter, r *http.Request) {
		output.HTML(w, http.StatusOK, "index", nil)
	})

	//admin
	ra := r.With(auth.Registry.GuardRequest("/denied", auth.AdminUser))
	dir := http.Dir(Config.Server.AppPath)
	fserver := http.StripPrefix("/admin/", http.FileServer(dir))
	ra.Get("/admin/*", func(w http.ResponseWriter, r *http.Request) {
		fserver.ServeHTTP(w, r)
	})

	log.Printf("starting on port %s", Config.Server.Port)
	http.ListenAndServe(Config.Server.Port, r)
}
