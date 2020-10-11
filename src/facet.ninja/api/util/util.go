package util

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func GenerateBase64UUID() string {
	uuid, _ := uuid.NewRandom()
	base64String := base64.StdEncoding.EncodeToString([]byte(uuid.String()))
	return base64String
}
