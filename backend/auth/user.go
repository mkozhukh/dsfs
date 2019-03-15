package auth

import (
	"log"

	"github.com/mkozhukh/roles"

	"github.com/mkozhukh/dsfs/backend/store"
)

type CurrentUser struct {
	IsOwner bool          `json:"owner"`
	Email   string        `json:"email"`
	Name    string        `json:"name"`
	Rights  []roles.Right `json:"rights"`
}

func (c *CurrentUser) Load(id uint) {
	// get full user info from DB
	user := store.Users.Get(id)

	c.Email = user.Email
	c.Name = user.Name
	c.IsOwner = user.IsOwner
	c.Rights = Registry.GetRights(roles.Role(id))
}

func CreateOwner(name, email string) *store.User {
	user := store.User{
		Email:   email,
		Name:    name,
		IsOwner: true,
		Rights:  roles.SerializeRights(AdminUser, UploadFiles),
	}

	store.Users.Save(&user)
	ReloadRoles()

	return &user
}

func ReloadRoles() {
	Registry.Reset()
	for _, user := range store.Users.GetAll() {
		log.Printf("%+v", user)
		data := roles.ParseRights(user.Rights)
		Registry.RegisterRole(roles.Role(user.ID), data...)
	}
}
