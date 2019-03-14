package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mkozhukh/go-remote"
	"github.com/mkozhukh/roles"

	"github.com/mkozhukh/dsfs/backend/store"
)

const (
	//CanAdminUser allows to manage users
	CanAdminUser roles.Right = iota
	//CanUploadFiles allows to store files
	CanUploadFiles
)

var auth = roles.NewRegistry()
var users *store.UserCollection

type CurrentUser struct {
	IsOwner bool          `json:"owner"`
	Email   string        `json:"email"`
	Name    string        `json:"name"`
	Rights  []roles.Right `json:"rights"`
}

func (c *CurrentUser) Load(id uint) {
	// get full user info from DB
	user := users.Get(id)

	c.Email = user.Email
	c.Name = user.Name
	c.IsOwner = user.Email == Config.Owner
	c.Rights = auth.GetRights(roles.Role(id))
}

func createDefaultUser(name, email string) *store.User {
	user := store.User{
		Email:  email,
		Name:   name,
		Rights: roles.SerializeRights(CanAdminUser, CanUploadFiles),
	}

	users.Save(&user)
	LoadRolesFromDB()

	return &user
}

func LoadRolesFromDB() {
	for _, user := range users.GetAll() {
		data := roles.ParseRights(user.Rights)
		auth.RegisterRole(roles.Role(user.ID), data...)
	}
}

func initUser(r chi.Router, s *remote.Server) {
	users, _ = store.NewUserCollection("./users.json")

	// providers for API
	s.RegisterProvider(func(r *http.Request) *roles.User {
		return roles.UserFromContext(r.Context())
	})
	s.RegisterProvider(func(r *http.Request) *CurrentUser {
		user := CurrentUser{}

		u := roles.UserFromContext(r.Context())
		if u.ID != 0 {
			user.Load(u.ID)
		}

		return &user
	})

	// static info for API
	s.RegisterConstant("rights", map[string]roles.Right{
		"CanUploadFiles": CanUploadFiles,
		"CanAdminUser":   CanAdminUser,
	})
	s.RegisterVariable("user", &CurrentUser{})
}
