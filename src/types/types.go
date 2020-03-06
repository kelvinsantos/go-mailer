package types

import (
	"github.com/aws/aws-sdk-go/service/ses"
)

type EmailInput struct {
	Error				error
	SendEmailInput		*ses.SendEmailInput
	SendRawEmailInput	*ses.SendRawEmailInput
}

type ApiResponse struct {
	Message		string		`json:"message"`
	Data		interface{}	`json:"data"`
	Success		bool		`json:"success"`
}

type SendEmailRequestJson struct {
	IsMockEmail		bool		`valid:"-"`
	Sender			string		`valid:"-"`
	ReplyTo			string		`valid:"-"`
	//ReplyTo		[]string	`valid:"-"`
	Subject			string		`valid:"-"`
	HtmlBody		string		`valid:"-"`
	ToRecipients	[]string	`valid:"-"`
	CcRecipients	[]string	`valid:"-"`
	BccRecipients	[]string	`valid:"-"`
	Attachments		[]Attachment	`valid:"-"`
}

type Attachment struct {
	Filename		string		`valid:"-"`
	ContentType		string		`valid:"-"`
	Content			string		`valid:"-"`
}