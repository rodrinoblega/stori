package email_service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

func CreateSESSession() *ses.SES {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // SES region
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	svc := ses.New(sess)

	return svc
}
