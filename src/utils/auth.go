package utils

import (
	"fmt"
	"log"
)

func ValidateAuthSecret(token string) string {
	auth := GetEnvWithDefault("GO_MAILER_AUTH_SECRET", "")

	if auth == "" {
		return "Auth Secret environment variable is missing"
	}

	if token == "" || token == "undefined" {
		log.Println("Authorization token empty")
		return "Authorization token empty"
	}

	if token == auth {
		fmt.Println("Validated auth secret: ", token)
		return ""
	}

	log.Println("Auth secret is invalid")
	return "Auth secret is invalid"
}
