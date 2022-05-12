package entity

type User struct {
	Id int64 `db:"id"`
	Nickname string `db:"nickname"`
	Email string `db:"email"`
	Balance uint64 `db:"balance"`
}

func NewUser(nickname, email string) *User {
	return &User{
		Nickname: nickname,
		Email: email,
		Balance: 0,
	}
}