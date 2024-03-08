package message

import "errors"

var (
	ErrFailedToSendMessage   = errors.New("message: failed to send message")
	ErrFailedToDeleteMessage = errors.New("message: failed to delete message")
	// ErrCreatePrivateChat is returned when cannot create private chat
	ErrCreatePrivateChat = errors.New("Cannot create private chat make sure you allow DM from server members")
	// ErrSendPrivateChat is returned when cannot send private chat
	ErrSendPrivateChat = errors.New("Cannot send private chat make sure you allow DM from server members")
)
