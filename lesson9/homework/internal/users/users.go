package users

type User struct {
	ID       int64
	Nickname string `validate:"min:1;max:30"`
	Email    string `validate:"min:1;max:30"`
}
