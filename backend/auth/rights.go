package auth

import "github.com/mkozhukh/roles"

const (
	//AdminUser allows to manage users
	AdminUser roles.Right = iota + 1
	//UploadFiles allows to store files
	UploadFiles
)

// All return all available rights
func All() map[string]roles.Right {
	return map[string]roles.Right{
		"UploadFiles": UploadFiles,
		"AdminUser":   AdminUser,
	}
}

//Registry is the default roles registry
var Registry *roles.Registry

func init() {
	Registry = roles.NewRegistry()
}
