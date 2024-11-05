package logic

import "fmt"

type FriendServiceError struct {
	Message string
}

func (e *FriendServiceError) Error() string {
	return fmt.Sprintf("friend service error: %s", e.Message)
}

var (
	ErrAlreadyFriends           = &FriendServiceError{Message: "already friends"}
	ErrCannotAddYourself        = &FriendServiceError{Message: "cannot add yourself as a friend"}
	ErrUserIdRequired           = &FriendServiceError{Message: "user id is required"}
	ErrUserIdOrFriendIdRequired = &FriendServiceError{Message: "user id and friend id are required"}
)
