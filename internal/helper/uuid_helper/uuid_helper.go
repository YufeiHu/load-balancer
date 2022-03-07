package uuid_helper

import "github.com/google/uuid"

func GenerateUuid() string {
	id := uuid.New()
	return id.String()
}
