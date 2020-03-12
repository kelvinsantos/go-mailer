package types

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type GetInboxRequestJson struct {
	Email			string
	Skip			int64
	Limit			int64
	Sort			int64
}

type GetMessageRequestJson struct {
	Email			string
	MessageId		string
}

type PartialMail struct {
	Id				string					`valid:"-" bson:"_id" json:"id"`
	Subject			string					`valid:"-" bson:"subject" json:"subject"`
	CreatedAt		primitive.DateTime		`valid:"-" bson:"created_at" json:"created_at"`
}

type Mail struct {
	Id				string					`valid:"-" bson:"_id" json:"id"`
	IsMockEmail		bool					`valid:"-" bson:"is_mock_email" json:"is_mock_email"`
	Sender			string					`valid:"-" bson:"sender" json:"sender"`
	ReplyTo			string					`valid:"-" bson:"reply_to" json:"reply_to"`
	Subject			string					`valid:"-" bson:"subject" json:"subject"`
	HtmlBody		string					`valid:"-" bson:"html_body" json:"html_body"`
	ToRecipients	[]string				`valid:"-" bson:"to_recipients" json:"to_recipients"`
	CcRecipients	[]string				`valid:"-" bson:"cc_recipients" json:"cc_recipients"`
	BccRecipients	[]string				`valid:"-" bson:"bcc_recipients" json:"bcc_recipients"`
	Attachments		[]Attachment			`valid:"-" bson:"attachments" json:"attachments"`
	CreatedAt		primitive.DateTime		`valid:"-" bson:"created_at" json:"created_at"`
}

type Attachment struct {
	Filename		string		`valid:"-"`
	ContentType		string		`valid:"-"`
	Content			string		`valid:"-"`
}

type AuditLog struct {
	Entry_Type		string
	ActionBy		ActionBy
	Company			Company
	Text			string
	Action			string
	OldValue		map[string]string
	NewValue		map[string]string
}

type ActionBy struct {
	Id				string
	FirstName		string
	LastName		string
	Email			string
}

type Company struct {
	Id				string
	Name			string
	Uen				string
}