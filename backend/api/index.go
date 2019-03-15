package api

import (
	"net/http"

	"github.com/mkozhukh/dsfs/backend/auth"
	"github.com/mkozhukh/go-remote"
	"github.com/mkozhukh/roles"
)

//New creates new remote API handler
func New() *remote.Server {
	api := remote.NewServer()

	api.RegisterWithGuard("admin", &UsersAPI{}, auth.Registry.CheckRequest(auth.AdminUser))

	// providers for API
	api.RegisterProvider(func(r *http.Request) *roles.User {
		return roles.UserFromContext(r.Context())
	})
	api.RegisterProvider(func(r *http.Request) *auth.CurrentUser {
		user := auth.CurrentUser{}

		u := roles.UserFromContext(r.Context())
		if u.ID != 0 {
			user.Load(u.ID)
		}

		return &user
	})

	// static info for API
	api.RegisterConstant("rights", auth.All())
	api.RegisterVariable("user", &auth.CurrentUser{})
	return api
}
