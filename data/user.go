package data

type (
	// User
	User struct {
		Token string `json:"-"`
		Id    string `json:"-"`
		Name  string `json:"username"`
		Pw    string `json:"-"`
	}
)

const (
	USR_COLLECTION = "user"

	USR_ID_NOTSET   = "The User.Id is required but hasn't been set"
	USR_NAME_NOTSET = "The User.Name is required but hasn't been set"
)

func (u *User) CanLogin() bool {
	return u.Name != "" && u.Pw != ""
}

func (u *User) IsLoggedIn() bool {
	return u.Token != ""
}

func (u *User) IsSet() bool {
	return u.Id != ""
}
