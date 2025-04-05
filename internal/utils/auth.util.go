package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("erro ao gerar API key: " + err.Error())
	}
	apiKey := base64.URLEncoding.EncodeToString(bytes)
	return apiKey, nil
}
