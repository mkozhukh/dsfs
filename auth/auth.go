package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"mkozhukh/dsfs/config"

	"github.com/alexedwards/scs/session"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"

	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
)

//AccessType is an enumeration of possible access levels
type AccessType int

//AccessContextType used as uniq identifier for the context property
type AccessContextType int

const (
	//NoneAccess means access denied
	NoneAccess AccessType = iota
	//WriteAccess gives full access
	WriteAccess
)

const (
	//AccessContext is a context property, which stores current access level
	AccessContext AccessContextType = iota
)

func getAccessByEmail(email string) AccessType {
	users := config.Config.Users
	for _, el := range users {
		if el == email {
			return WriteAccess
		}
	}

	return NoneAccess
}

//AddRoutes attach handlers to auth routes
func AddRoutes(format *render.Render) http.Handler {
	r := chi.NewRouter()
	config := &config.Config

	//register auth providers
	url := fmt.Sprintf("%s/auth/gplus/callback", config.Google.Callback)
	goth.UseProviders(
		gplus.New(config.Google.Key, config.Google.Secret, url),
	)

	//add routes
	r.Get("/{provider}/callback", callback)
	r.Get("/{provider}/login", login)
	r.Get("/{provider}/logout", logout)
	r.Get("/denied", func(res http.ResponseWriter, req *http.Request) {
		format.HTML(res, http.StatusOK, "auth/denied", nil)
	})

	return r
}

func callback(res http.ResponseWriter, req *http.Request) {
	user, err := completeUserAuth(res, req)
	if err != nil {
		log.Error("Can't complete user's authentication")
		log.Error(err)
	}
	gateway(res, req, user.Email)
}

func login(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	if user, err := completeUserAuth(res, req); err == nil {
		gateway(res, req, user.Email)
	} else {
		beginAuthHandler(res, req)
	}
}

func logout(res http.ResponseWriter, req *http.Request) {
	logoutHandler(res, req)
	redirect(res, "/")
}

//GetAccess returns acess type for the current user
func GetAccess(req *http.Request) AccessType {
	email, err := session.GetString(req, "email")
	if err != nil {
		return NoneAccess
	}

	return getAccessByEmail(email)
}

//ForceLogin is a middleware that check login status
func ForceLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access := GetAccess(r)
		if access == NoneAccess &&
			!strings.HasPrefix(r.URL.Path, "/auth/") &&
			!strings.HasPrefix(r.URL.Path, "/debug/") &&
			!strings.HasPrefix(r.URL.Path, "/"+config.Config.Path+"/") {
			http.Redirect(w, r, "/auth/gplus/login", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), AccessContext, access)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func gateway(res http.ResponseWriter, req *http.Request, email string) {
	access := getAccessByEmail(email)
	if access == NoneAccess {
		redirect(res, "/auth/denied")
		return
	}

	if session.PutString(req, "email", email) != nil {
		log.Error("Can't save auth info into session")
		redirect(res, "/auth/denied")
		return
	}

	redirect(res, "/")
}

func redirect(res http.ResponseWriter, url string) {
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
