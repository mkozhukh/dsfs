package store

type User struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Rights  string `json:"rights"`
	IsOwner bool   `json:"owner"`
}

var Users *UserCollection

func init() {
	Users, _ = NewUserCollection("./data/users.json")
}
