package api

import (
	"errors"

	"github.com/mkozhukh/dsfs/backend/auth"
	"github.com/mkozhukh/dsfs/backend/store"
	"github.com/mkozhukh/roles"
)

type UsersAPI struct{}

// GetUsers returns all available Users
func (a *UsersAPI) GetUsers(role *roles.User) []store.User {
	role.Guard(auth.AdminUser)

	return store.Users.GetAll()
}

// SaveUser is a CRUD handler for Users collection
func (a *UsersAPI) SaveUser(id uint, action string, s store.User, role *roles.User) (*store.User, error) {
	role.Guard(auth.AdminUser)

	if s.ID != 0 {
		//check that we not
		old := store.Users.Get(s.ID)

		// only owner can edit his data
		if old.IsOwner && role.ID != old.ID {
			return nil, errors.New("Acces denied")
		}

		s.IsOwner = old.IsOwner
	}

	if action == "delete" {
		store.Users.Delete(id)
		auth.ReloadRoles()
		return nil, nil
	}

	similar := store.Users.First(func(u *store.User) bool {
		return u.Email == s.Email
	})
	// prevent taking email of other user
	if similar.ID != 0 && similar.ID != s.ID {
		return nil, errors.New("Acces denied")
	}

	store.Users.Save(&s)
	auth.ReloadRoles()
	return &s, nil
}
