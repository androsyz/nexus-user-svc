package constant

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrForbiddenAccess    = errors.New("forbidden access")
	ErrMissingCredentials = errors.New("email and password are required")
	ErrHashingPassword    = errors.New("failed to hash password")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidClaims       = errors.New("invalid claims in token")
	ErrInvalidSubject      = errors.New("invalid subject in token")
	ErrGeneratingJWT       = errors.New("failed to generate token")
)

var (
	ErrMsgMarshal   = "failed to marshal message"
	ErrMsgUnmarshal = "failed to unmarshal message"
	ErrMsgPublish   = "failed to publish message to redis"
	ErrMsgSubscribe = "failed to subscribe message from redis"
	ErrMsgSubsFull  = "subscriber channel full, skipping message"
)

func ErrMissingField(field string) error {
	return fmt.Errorf("%s is required", field)
}

func ErrCreatingField(field string) error {
	return fmt.Errorf("failed to create %s", field)
}

func ErrGetField(field string) error {
	return fmt.Errorf("failed to get %s", field)
}

func ErrWithMsg(errMsg, err error) error {
	return fmt.Errorf("%w: %v", errMsg, err)
}
