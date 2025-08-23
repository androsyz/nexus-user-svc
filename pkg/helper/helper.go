package helper

import (
	"github.com/segmentio/ksuid"
)

func NewKSUID() string {
	return ksuid.New().String()
}

func ParseKSUID(id string) (ksuid.KSUID, error) {
	return ksuid.Parse(id)
}
