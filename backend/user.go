package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mkozhukh/go-remote"
	"github.com/mkozhukh/roles"
)

const (
	//CanAdminUser allows to manage users
	CanAdminUser roles.Right = iota
	//CanEditData allows to auth app and change data
	CanUploadFiles
)

var auth = roles.NewRegistry()

type CurrentUser struct {
	IsOwner bool          `json:"owner"`
	Email   string        `json:"email"`
	Name    string        `json:"name"`
	Rights  []roles.Right `json:"rights"`
}

func (c *CurrentUser) Load(id uint) {
	// get full user info from DB
	user := User{}
	db.Find(&user, id)

	c.Email = user.Email
	c.Name = user.Name
	c.IsOwner = user.Email == Config.Owner
	c.Rights = auth.GetRights(roles.Role(id))
}

func createDefaultUser(name, email string) *User {
	user := User{
		Email:  email,
		Name:   name,
		Rights: roles.SerializeRights(CanAdminUser, CanUploadFiles),
	}

	db.Save(&user)
	Config.LoadRoles()

	return &user
}

func (c *AppConfig) LoadRolesFromDB() {
	users := make([]User, 0)
	db.Find(&users)

	for _, user := range users {
		data := roles.ParseRights(user.Rights)
		auth.RegisterRole(roles.Role(user.ID), data...)
	}
}

func initUser(r chi.Router, s *remote.Server) {
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
