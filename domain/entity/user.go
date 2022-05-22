package entity

type User struct {
	Id int64 `db:"id"`
	Nickname string `db:"nickname"`
	Email string `db:"email"`
	Balance uint64 `db:"balance"`
	Password string `db:"password"`
}

func NewUser(nickname, email string, hashedPW string) *User {
	return &User{
		Nickname: nickname,
		Email: email,
		Balance: 0,
		Password: hashedPW,
	}
}

