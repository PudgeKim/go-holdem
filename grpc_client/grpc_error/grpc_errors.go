package grpc_error

import "errors"

var (
	UserNotExist = errors.New("user does not exist")
	InvalidUser  = errors.New("user id does not exist")
)
