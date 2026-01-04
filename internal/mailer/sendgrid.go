package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, toEmail string, data any, isSandbox bool) error {
	// Implementation for sending email using SendGrid
	from := mail.NewEmail(FromName,s.fromEmail)
	to := mail.NewEmail(username, toEmail) 
	
	//template parsing and building
	tmpl,err:=template.ParseFS(FS,"templates/"+templateFile)
	if err!=nil{
		return err
	}
	subject := new(bytes.Buffer)
	err= tmpl.ExecuteTemplate(subject,"subject",data)
	if err!=nil{
		return err
	}
	body := new(bytes.Buffer)
	err= tmpl.ExecuteTemplate(body,"body",data)
	if err!=nil{
		return err
	}
	message := mail.NewSingleEmail(from,subject.String(),to,"",body.String());
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})
	for i:=0;i<maxRetires;i++{
		response,err := s.client.Send(message)
		if err!=nil{
			log.Printf("sendgrid send email attempt %d failed: %v",i+1,err)
			//exponential backoff
			time.Sleep(time.Second*time.Duration(i+1))
			continue;
		}
		log.Printf("sendgrid send email response: %d",response.StatusCode)
		return nil;
	}
	return fmt.Errorf("failed to send email after max retries.")
}
