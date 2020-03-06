package store

import (
	"MailerGo/src/decorator"
	"MailerGo/src/env"
	"MailerGo/src/types"
	"MailerGo/src/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// Replace customer@gmail.com with your "From" address.
	// This address must be verified with Amazon SES.
	// Sender = "customer@gmail.com"

	// Replace s.kelvinjohn@gmail.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	// Recipient = "s.kelvinjohn@gmail.com"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	// Subject = "Amazon SES Test (AWS SDK for Go)"

	// The HTML body for the email.
	// HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
	// 	"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
	// 	"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	// The email body for recipients with non-HTML email clients.
	// TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func buildEmailInput(requestJson types.SendEmailRequestJson) types.EmailInput {
	var emailInput types.EmailInput

	//Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  []*string{},
			CcAddresses:  []*string{},
			BccAddresses: []*string{},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(requestJson.HtmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(requestJson.Subject),
			},
		},
		ReplyToAddresses: []*string{
			aws.String(requestJson.ReplyTo),
		},
		Source: aws.String(requestJson.Sender),
	}

	//append all to recipients
	if len(requestJson.ToRecipients) > 0 {
		for _, to := range requestJson.ToRecipients {
			input.Destination.ToAddresses = append(input.Destination.ToAddresses, aws.String(to))
		}
	}

	// append all cc recipients
	if len(requestJson.CcRecipients) > 0 {
		for _, cc := range requestJson.CcRecipients {
			input.Destination.CcAddresses = append(input.Destination.CcAddresses, aws.String(cc))
		}
	}

	// append all bcc recipients
	if len(requestJson.BccRecipients) > 0 {
		for _, bcc := range requestJson.BccRecipients {
			input.Destination.BccAddresses = append(input.Destination.BccAddresses, aws.String(bcc))
		}
	}

	emailInput.SendEmailInput = input

	return emailInput
}

func buildRawEmailInput(requestJson types.SendEmailRequestJson) types.EmailInput {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	var emailInput types.EmailInput

	// email main header:
	header := make(textproto.MIMEHeader)
	header.Set("From", requestJson.Sender)
	header.Set("Reply-To", requestJson.ReplyTo)
	header.Set("Return-Path", requestJson.Sender)
	header.Set("Subject", requestJson.Subject)
	header.Set("Content-Language", "en-US")
	header.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	header.Set("MIME-Version", "1.0")

	// todo set custom headers here
	header.Set("X-SES-CONFIGURATION-SET", "Platform-dev")
	header.Set("List-Unsubscribe", "<mailto:customer@gmail.com?subject=unsubscribe>")

	_, err := writer.CreatePart(header)
	if err != nil {
		emailInput.Error = err
	}

	// set to recipients in header
	var csvToRecipients string
	toRecipients := requestJson.ToRecipients
	if len(toRecipients) > 0 {
		for index, to := range toRecipients {
			csvToRecipients += to
			if index != len(toRecipients) -1 {
				csvToRecipients += ","
			}
		}
		header.Set("To", csvToRecipients)
	}

	// set cc recipients in header
	var csvCcRecipients string
	ccRecipients := requestJson.CcRecipients
	if len(ccRecipients) > 0 {
		for index, to := range ccRecipients {
			csvCcRecipients += to
			if index != len(ccRecipients) -1 {
				csvCcRecipients += ","
			}
		}
		header.Set("Cc", csvCcRecipients)
	}

	// set bcc recipients in header
	var csvBccRecipients string
	bccRecipients := requestJson.BccRecipients
	if len(bccRecipients) > 0 {
		for index, to := range bccRecipients {
			csvBccRecipients += to
			if index != len(bccRecipients) -1 {
				csvBccRecipients += ","
			}
		}
		header.Set("Bcc", csvBccRecipients)
	}

	// body:
	header = make(textproto.MIMEHeader)
	header.Set("Content-Transfer-Encoding", "7bit")
	header.Set("Content-Type", "text/html; charset=UTF-8")
	part, err := writer.CreatePart(header)
	if err != nil {
		emailInput.Error = err
	}

	_, err = part.Write([]byte(requestJson.HtmlBody))
	if err != nil {
		emailInput.Error = err
	}

	files := requestJson.Attachments
	if files != nil {
		for _, file := range files {
			header = make(textproto.MIMEHeader)
			header.Set("Content-Type", file.ContentType + ";name=\""+file.Filename+"\"")
			header.Set("Content-Description", file.Filename)
			header.Set("Content-Disposition", "attachment;filename=\""+file.Filename+"\"")
			header.Set("Content-Transfer-Encoding", "base64")

			data, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				log.Fatal("error:", err)
			}
			//fmt.Printf("%q\n", []byte(string(data)))

			part, err = writer.CreatePart(header)
			if err != nil {
				emailInput.Error = err
			}

			_, err = part.Write(data)
			if err != nil {
				emailInput.Error = err
			}

			err = writer.Close()
			if err != nil {
				emailInput.Error = err
			}
		}
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		emailInput.Error = fmt.Errorf("invalid e-mail content")
	}

	log.Println(s)

	s = strings.SplitN(s, "\n", 2)[1]

	//log.Println(s)

	raw := ses.RawMessage{
		Data: []byte(s),
	}

	input := &ses.SendRawEmailInput{
		Destinations: []*string{},
		Source:       aws.String(requestJson.Sender),
		RawMessage:   &raw,
	}

	if len(toRecipients) > 0 {
		for _, to := range toRecipients {
			input.Destinations = append(input.Destinations, aws.String(to))
		}
	}

	emailInput.SendRawEmailInput = input

	return emailInput
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var requestJson types.SendEmailRequestJson

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &requestJson)

	buildEmailInput := decorator.BuildEmailInputDecorate(buildEmailInput)
	input := buildEmailInput(requestJson)

	if input.Error != nil {
		utils.HandleError(w, input.Error)
	}

	svc, err := createSesSession(input)
	if err != nil {
		utils.HandleError(w, err)
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input.SendEmailInput)

	// Display error messages if they occur.
	if err != nil {
		utils.HandleSesError(w, err)
	}

	fmt.Println("Email Sent to address: ")
	fmt.Println(requestJson.ToRecipients)
	fmt.Println(result)

	return
}

func SendRawEmail(w http.ResponseWriter, r *http.Request) {
	var requestJson types.SendEmailRequestJson

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &requestJson)

	buildRawEmailInput := decorator.BuildEmailInputDecorate(buildRawEmailInput)
	input := buildRawEmailInput(requestJson)

	if input.Error != nil {
		utils.HandleError(w, input.Error)
	}

	svc, err := createSesSession(input)
	if err != nil {
		utils.HandleError(w, err)
	}

	// Attempt to send the email.
	result, err := svc.SendRawEmail(input.SendRawEmailInput)

	// Display error messages if they occur.
	if err != nil {
		utils.HandleSesError(w, err)
	}

	fmt.Println("Email Sent to address: ")
	fmt.Println(requestJson.ToRecipients)
	fmt.Println(result)

	return
}

func createSesSession(input types.EmailInput) (*ses.SES, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(env.GO_MAILER_AWS_REGION),
		Credentials: credentials.NewStaticCredentials(env.GO_MAILER_AWS_ACCESS_KEY, env.GO_MAILER_AWS_SECRET_KEY, ""),
	})

	if err != nil {
		return nil, err
	}

	// Create an SES session.
	svc := ses.New(sess)

	return svc, nil
}