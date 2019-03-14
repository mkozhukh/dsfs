package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"
	"github.com/markbates/goth/providers/google"
	"github.com/mkozhukh/login"
	"github.com/mkozhukh/roles"

	"github.com/mkozhukh/dsfs/backend/store"
)

type loginHandler struct {
	session *scs.Manager
}

func (l loginHandler) Login(req *http.Request, res http.ResponseWriter, email string) string {
	user := users.First(func(u *store.User) bool {
		return u.Email == email
	})

	// unknown user
	if user.ID == 0 {
		if email == Config.Owner {
			user = createDefaultUser("Admin", Config.Owner)
		}
		return "/denied"
	}

	if l.session.Load(req).PutInt(res, "userID", int(user.ID)) != nil {
		log.Print("Can't save auth info into session")
		return "/denied"
	}

	return "/admin"
}

func (l loginHandler) Logout(req *http.Request, res http.ResponseWriter) string {
	if l.session.Load(req).Remove(res, "userID") != nil {
		log.Print("Can't delete auth info from session")
		return "/admin"
	}

	return "/"
}

func initLogin(r chi.Router, session *scs.Manager) {
	// store user object in context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, _ := session.Load(r).GetInt("userID")
			user := auth.NewUser(uint(userID), roles.Role(userID))
			ctx := roles.UserToContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	//auth
	login.SetSession(session)

	// add login|logout routes
	subRoute := chi.NewRouter()
	login.SetProvider(
		google.New(Config.Google.Key, Config.Google.Secret, Config.Google.Callback),
		subRoute,
		"/login", "/logout", "/callback",
		loginHandler{session},
	)
	r.Mount("/auth", subRoute)
}
