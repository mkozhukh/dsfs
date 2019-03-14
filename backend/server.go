// go:generate go-bindata static/

package main

import (
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

	// Routes
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	api := remote.NewServer()

	initDB()
	initAPI(api)

	initLogin(r, session)
	initUser(r, api)
	initUpload(r, output)

	//static assets
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		output.HTML(w, http.StatusOK, "index", nil)
	})

	http.ListenAndServe(":"+Config.Server.Port, r)
}
