package email

import (
	"bytes"
	"context"
	"fmt"
	htemplate "html/template"
	"mime/multipart"
	"net/smtp"
	"sort"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
)

type Sender interface {
	Send(ctx context.Context, values map[string]string, files map[string]*multipart.FileHeader) error
}

func NewSMTPSender(cfg SMTPConfig) (*SMTPSender, error) {
	s := SMTPSender{
		cfg: cfg,
	}

	var err error

	if strings.ContainsAny(cfg.Subject, "{}") {
		s.subjectTmpl, err = template.New("subject").Parse(cfg.Subject)
		if err != nil {
			return nil, fmt.Errorf("failed to parse subject template: %w", err)
		}
	}

	if cfg.BodyTemplateFile != "" {
		s.bodyTmpl, err = template.ParseFiles(cfg.BodyTemplateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse body template file: %w", err)
		}
	}

	if cfg.BodyHTMLTemplateFile != "" {
		s.bodyHTMLTmpl, err = htemplate.ParseFiles(cfg.BodyHTMLTemplateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse body html template file: %w", err)
		}
	}

	return &s, nil
}

type SMTPSender struct {
	cfg          SMTPConfig
	subjectTmpl  *template.Template
	bodyTmpl     *template.Template
	bodyHTMLTmpl *htemplate.Template
}

// SMTPConfig defines SMTP configuration.
type SMTPConfig struct {
	From                 string `split_words:"true"`
	Password             string `split_words:"true"`
	To                   string `split_words:"true"`
	Host                 string `split_words:"true"`
	Port                 int    `split_words:"true" default:"587"`
	Subject              string `split_words:"true" default:"New Form Received"`
	ReplyTo              string `split_words:"true"`
	BodyTemplateFile     string `split_words:"true"`
	BodyHTMLTemplateFile string `split_words:"true"`
}

func (s *SMTPSender) EmailSender() Sender {
	return s
}

func (s *SMTPSender) Send(_ context.Context, values map[string]string, files map[string]*multipart.FileHeader) error {
	smtpHost := s.cfg.Host
	smtpPort := s.cfg.Port
	from := s.cfg.From
	to := []string{s.cfg.To}
	password := s.cfg.Password
	auth := smtp.PlainAuth("", from, password, smtpHost)

	e := email.NewEmail()
	e.From = from
	e.To = to

	if s.cfg.ReplyTo != "" {
		e.ReplyTo = []string{values[s.cfg.ReplyTo]}
	}

	if s.subjectTmpl != nil {
		buf := bytes.NewBuffer(nil)
		if err := s.subjectTmpl.Execute(buf, values); err != nil {
			return err
		}
		e.Subject = buf.String()
	} else {
		e.Subject = s.cfg.Subject
	}

	if s.bodyTmpl != nil {
		buf := bytes.NewBuffer(nil)
		if err := s.bodyTmpl.Execute(buf, values); err != nil {
			return err
		}
		e.Text = buf.Bytes()
	} else {
		text := ""

		keys := make([]string, 0, len(values))
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			text += strings.Title(k) + ":\n" + values[k] + "\n\n"
		}

		e.Text = []byte(text)
	}

	if s.bodyHTMLTmpl != nil {
		buf := bytes.NewBuffer(nil)
		if err := s.bodyHTMLTmpl.Execute(buf, values); err != nil {
			return err
		}
		e.HTML = buf.Bytes()
	}

	for k, v := range files {
		f, err := v.Open()
		if err != nil {
			return err
		}
		_, err = e.Attach(f, k+"_"+v.Filename, v.Header.Get("Content-Type"))
		if err != nil {
			return err
		}
	}

	return e.Send(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth)
}
