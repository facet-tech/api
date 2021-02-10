package notification

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"strings"
)

type Notification struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Company     string `json:"company,omitempty"`
	Contact     string `json:"contact,required"`
	Subject     string `json:"subject,required"`
	Message     string `json:"message,omitempty"`
}

type SlackMessage struct {
	Message     string `json:"message"`
}

const (
	slackFormatBlock = "\n>"
	slackFormatBold  = "*"
	topicArn         = "arn:aws:sns:us-west-2:935571265336:customer-leads"
)

func slackBold(string string) string  {
	return slackFormatBold + string + slackFormatBold
}

func slackField(key string, value string) string  {
	if value == "" {
		return ""
	} else {
		return slackBold(key) + ": " + strings.ReplaceAll(value, "\n", slackFormatBlock) + slackFormatBlock
	}
}

func (notification *Notification) createSlackMessage() SlackMessage  {
	return SlackMessage{
		slackBold(notification.Subject) + slackFormatBlock +
				slackField("Company", notification.Company) +
				slackField("FirstName", notification.FirstName) +
			    slackField("LastName", notification.LastName) +
				slackField("Contact", notification.Contact) +
				slackField("Message", notification.Message)}
}

func (notification *Notification) SendBatch() error {
	message, error := json.Marshal(notification.createSlackMessage())
	if error == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		svc := sns.New(sess, aws.NewConfig().WithRegion("us-west-2"))
		_, error = svc.Publish(&sns.PublishInput{
			Message:  aws.String(string(message)),
			TopicArn: aws.String(topicArn),
		})
	}
	return error
}