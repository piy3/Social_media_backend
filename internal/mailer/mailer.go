package mailer

import "embed"

//go:embed templates 
var FS embed.FS

const (
	FromName = "GoSocial"
	maxRetires = 3
	UserWelcomeTemplate = "user_welcome.tmpl"
)

type Client interface {
	Send(templateFile,username,toEmail string,data any,isSandbox bool) error 
}