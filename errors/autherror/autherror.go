package autherror

import "errors"

var (
	InvalidToken = errors.New("invalid token")
	UserNotFound = errors.New("user does not exist")
)