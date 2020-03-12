package decorator

import (
	"github.com/asaskevich/govalidator"
	"kelvin.com/mailer/src/env"
	"kelvin.com/mailer/src/types"
	"strings"
)

type Object func(requestJson types.SendEmailRequestJson) types.EmailInput

func BuildEmailInputDecorate(fn Object) Object {
	return func(requestJson types.SendEmailRequestJson) types.EmailInput {
		//log.Println("Starting the execution with the requestJson", requestJson)

		var input types.EmailInput

		isRequestValidated, err := govalidator.ValidateStruct(requestJson)
		if err != nil {
			input.Error = err
		}

		if isRequestValidated {
			// check if isMockEmail is equal to true
			// if true send emails to disposable emails
			// if false send emails to the real email address
			if requestJson.IsMockEmail == true {
				var toRecipients []string
				for _, element := range requestJson.ToRecipients {
					toRecipients = append(toRecipients, strings.Replace(element, "@", "_", -1) + env.GO_MAILER_MOCK_EMAIL)
				}
				requestJson.ToRecipients = toRecipients

				var ccRecipients []string
				for _, element := range requestJson.CcRecipients {
					ccRecipients = append(ccRecipients, strings.Replace(element, "@", "_", -1) + env.GO_MAILER_MOCK_EMAIL)
				}
				requestJson.CcRecipients = ccRecipients

				var bccRecipients []string
				for _, element := range requestJson.BccRecipients {
					bccRecipients = append(bccRecipients, strings.Replace(element, "@", "_", -1) + env.GO_MAILER_MOCK_EMAIL)
				}
				requestJson.BccRecipients = bccRecipients
			}

			input = fn(requestJson)
		}

		//log.Println("Execution is completed with the input", input)

		return input
	}
}