package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"github.com/dubass83/go-micro-mailer/util"
	"github.com/vanng822/go-premailer/premailer"
	"github.com/wneessen/go-mail"
)

type EmailSender interface {
	SendEmail(email Message, errChan chan error)
}

type Message struct {
	From        string
	FromEmail   string
	Subject     string
	To          []string
	CC          []string
	BCC         []string
	Data        any
	Message     map[string]any
	AttachFiles []string
	Template    string
}

type MailTrapSender struct {
	From        string
	FromEmail   string
	Login       string
	Password    string
	SMTPHost    string
	SMTPPort    int
	SMTPAuth    mail.SMTPAuthType
	TemplateDir string
	Wg          *sync.WaitGroup
	MailerChan  chan Message
	ErrChan     chan error
	DoneChan    chan bool
}

func NewMailSender(conf util.Config) (EmailSender, error) {
	switch conf.EmailService {
	case "mailtrap":
		return &MailTrapSender{
			From:        conf.SenderName,
			FromEmail:   conf.SenderEmail,
			Login:       conf.EmailLogin,
			Password:    conf.EmailPassword,
			SMTPHost:    "sandbox.smtp.mailtrap.io",
			SMTPPort:    2525,
			SMTPAuth:    mail.SMTPAuthPlain,
			TemplateDir: conf.TemplateDir,
			Wg:          &sync.WaitGroup{},
		}, nil
	default:
		return nil, fmt.Errorf("not implemented any other mail service except mailtrap")
	}
}

func (sender *MailTrapSender) SendEmail(
	email Message,
	errChan chan error,
) {
	if email.Template == "" {
		email.Template = "mail"
	}
	if email.From == "" {
		email.From = sender.From
	}
	if email.FromEmail == "" {
		email.FromEmail = sender.FromEmail
	}

	m := mail.NewMsg()
	if err := m.FromFormat(email.From, email.FromEmail); err != nil {
		errChan <- fmt.Errorf("failed to set from address: %s", err)
	}
	if err := m.To(email.To...); err != nil {
		errChan <- fmt.Errorf("failed to set To address: %s", err)
	}
	if err := m.Cc(email.CC...); err != nil {
		errChan <- fmt.Errorf("failed to set CC address: %s", err)
	}
	if err := m.Bcc(email.BCC...); err != nil {
		errChan <- fmt.Errorf("failed to set BCC address: %s", err)
	}
	m.Subject(email.Subject)
	// generate and set to the message text plain body
	templPlain := fmt.Sprintf("%s/%s.plain.gohtml", sender.TemplateDir, email.Template)
	contentPlain, err := builPlainTextMessage(templPlain, email.Message)
	if err != nil {
		errChan <- fmt.Errorf("failed to generate plain text message: %s", err)
	}
	m.SetBodyString(mail.TypeTextPlain, contentPlain)
	// generate and set to the message alternative html formated body
	templFormated := fmt.Sprintf("%s/%s.html.gohtml", sender.TemplateDir, email.Template)
	contentHtml, err := buildHTMLMessage(templFormated, email.Message)
	if err != nil {
		errChan <- fmt.Errorf("failed to generate html formated message: %s", err)
	}
	m.AddAlternativeString(mail.TypeTextHTML, contentHtml)

	for _, file := range email.AttachFiles {
		m.AttachFile(file)
	}

	c, err := mail.NewClient(
		sender.SMTPHost,
		mail.WithPort(sender.SMTPPort),
		mail.WithSMTPAuth(sender.SMTPAuth),
		mail.WithUsername(sender.Login),
		mail.WithPassword(sender.Password),
	)
	if err != nil {
		errChan <- fmt.Errorf("failed to create mail client: %s", err)
	}

	if err = c.DialAndSend(m); err != nil {
		errChan <- err
	}
}

func buildHTMLMessage(templ string, message map[string]any) (string, error) {
	t, err := template.New("email-html").ParseFiles(templ)
	if err != nil {
		return "", fmt.Errorf("failed to create template from %s: %s", templ, err)
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", message); err != nil {
		return "", fmt.Errorf("failed execute template with message %v: %s", message, err)
	}

	formattedMessage, err := inlineCSS(tpl.String())
	if err != nil {
		return "", fmt.Errorf("failed generate inline CSS message from template: %s", err)
	}
	return formattedMessage, nil
}

func inlineCSS(fm string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(fm, &options)
	if err != nil {
		return "", fmt.Errorf("failed create premailer from string %s: %s", fm, err)
	}

	html, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("failed transform premailer to string: %s", err)
	}
	return html, nil
}

func builPlainTextMessage(templ string, message map[string]any) (string, error) {
	t, err := template.New("email-plain").ParseFiles(templ)

	if err != nil {
		return "", fmt.Errorf("failed to create template from %s: %s", templ, err)
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", message); err != nil {
		return "", fmt.Errorf("failed execute template with message %v: %s", message, err)
	}

	return tpl.String(), nil
}
