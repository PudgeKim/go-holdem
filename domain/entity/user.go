package entity

type User struct {
	Id int64 `db:"id"`
	Nickname string `db:"nickname"`
	Email string `db:"email"`
	Balance uint64 `db:"balance"`
	password string `db:"password"`
}

func NewUser(nickname, email string, hashedPW string) *User {
	return &User{
		Nickname: nickname,
		Email: email,
		Balance: 0,
		password: hashedPW,
	}
}

func (u User) GetPassword() string {
	return u.password
}