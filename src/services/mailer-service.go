package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mime/multipart"
	"net/textproto"
	"kelvin.com/mailer/src/env"
	"kelvin.com/mailer/src/types"
	"strings"
	"time"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// Replace s.kelvinjohn@gmail.com with your "From" address.
	// This address must be verified with Amazon SES.
	// Sender = "s.kelvinjohn@gmail.com"

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

func BuildRawEmailInput(requestJson types.SendEmailRequestJson) types.EmailInput {
	log.Println(">>> Sending e-mail with subject: " + requestJson.Subject)
	log.Printf("To: %s", requestJson.ToRecipients)
	log.Printf("Cc: %s", requestJson.CcRecipients)
	log.Printf("Bcc: %s", requestJson.BccRecipients)

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
	header.Set("List-Unsubscribe", "<mailto:s.kelvinjohn@gmail.com?subject=unsubscribe>")

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

	s = strings.SplitN(s, "\n", 2)[1]

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

func CreateSesSession(input types.EmailInput) (*ses.SES, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(env.GO_MAILER_AWS_REGION),
	})

	if env.GO_MAILER_AWS_ACCOUNT != "" {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(env.GO_MAILER_AWS_REGION),
			Credentials: credentials.NewSharedCredentials("", env.GO_MAILER_AWS_ACCOUNT),
		})
	}

	if err != nil {
		return nil, err
	}

	// Create an SES session.
	svc := ses.New(sess)

	return svc, nil
}

func GetInboxService(getInboxRequestJson types.GetInboxRequestJson) (error, []types.PartialMail) {
	log.Println(">>> Retrieving inbox for " + getInboxRequestJson.Email)

	collection := Client.Database(env.GO_MAILER_DB_NAME).Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var filterOptions bson.D
	if getInboxRequestJson.Email != "" {
		filterOptions = bson.D{{ Key: "to_recipients", Value: getInboxRequestJson.Email }}
	}

	findOptions := options.Find()
	findOptions.SetSkip(getInboxRequestJson.Skip)
	findOptions.SetLimit(getInboxRequestJson.Limit)
	findOptions.SetSort(bson.D{{"created_at", getInboxRequestJson.Sort }})

	var partialMail types.PartialMail
	cur, err := collection.Find(ctx, filterOptions, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	var partialMails []types.PartialMail
	for cur.Next(ctx) {
		err := cur.Decode(&partialMail)
		partialMails = append(partialMails, partialMail)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil, partialMails
}

func GetMessageService(getMessageRequestJson types.GetMessageRequestJson) (error, types.Mail) {
	log.Println(">>> Retrieving message " + getMessageRequestJson.MessageId)

	collection := Client.Database(env.GO_MAILER_DB_NAME).Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, _ := primitive.ObjectIDFromHex(getMessageRequestJson.MessageId)
	filterOptions := bson.D{
		{ "_id", objId },
		{ "to_recipients", getMessageRequestJson.Email },
	}

	log.Println(filterOptions)

	var mail types.Mail
	cur, err := collection.Find(ctx, filterOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		err := cur.Decode(&mail)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil, mail
}

func SaveEmail(requestJson types.SendEmailRequestJson){
	//log.Println(">>> Saving e-mail with subject: " + requestJson.Subject)
	//log.Printf("To: %s", requestJson.ToRecipients)
	//log.Printf("Cc: %s", requestJson.CcRecipients)
	//log.Printf("Bcc: %s", requestJson.BccRecipients)

	collection := Client.Database(env.GO_MAILER_DB_NAME).Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := collection.InsertOne(ctx, bson.D{
		{ "attachments", requestJson.Attachments },
		{ "reply_to", requestJson.ReplyTo },
		{ "sender", requestJson.Sender },
		{ "html_body", requestJson.HtmlBody },
		{ "subject", requestJson.Subject },
		{ "bcc_recipients", requestJson.BccRecipients },
		{ "cc_recipients", requestJson.CcRecipients },
		{ "to_recipients", requestJson.ToRecipients },
		{ "is_mock_email", requestJson.IsMockEmail },
		{ "created_at", time.Now() },
	})

	if err != nil {
		log.Fatal(err)
	}
}