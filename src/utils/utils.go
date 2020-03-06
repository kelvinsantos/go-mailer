package utils

import (
	"MailerGo/src/env"
	"MailerGo/src/types"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Init() {
	fmt.Println("Loading env...")

	err := godotenv.Load("mailer.env")

	if err != nil {
		log.Println("Error loading mailer.env file ", err)
	}

	env.GO_MAILER_AUTH_SECRET = GetEnvWithDefault("GO_MAILER_AUTH_SECRET", "")
	env.GO_MAILER_AWS_ACCESS_KEY = GetEnvWithDefault("GO_MAILER_AWS_ACCESS_KEY", "")
	env.GO_MAILER_AWS_SECRET_KEY = GetEnvWithDefault("GO_MAILER_AWS_SECRET_KEY", "")
	env.GO_MAILER_AWS_REGION = GetEnvWithDefault("GO_MAILER_AWS_REGION", "us-west-2")
	env.GO_MAILER_PORT = GetEnvWithDefault("GO_MAILER_PORT", "")
	env.GO_MOCK_EMAIL = GetEnvWithDefault("GO_MOCK_EMAIL", "")

	log.Println("Succesfully loaded environment variables")
}

func GetEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func HandleSuccess(w http.ResponseWriter, message string, data interface{}) {
	json.NewEncoder(w).Encode(types.ApiResponse{
		Message: message,
		Data:    data,
		Success: false,
	})
	return
}

func HandleError(w http.ResponseWriter, err error) {
	fmt.Println("Generic error: " + err.Error())
	json.NewEncoder(w).Encode(types.ApiResponse{
		Message: err.Error(),
		Success: false,
	})
	return
}

func HandleSesError(w http.ResponseWriter, err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
	json.NewEncoder(w).Encode(types.ApiResponse{
		Message: err.Error(),
		Success: false,
	})
	return
}

func EncodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}

func DecodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	base64.StdEncoding.Decode(base64Text, []byte(message))
	fmt.Printf("base64: %s\n", base64Text)
	return string(base64Text)
}
