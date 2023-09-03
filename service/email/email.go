package email

import (
	"net/url"
	"strconv"

	"github.com/gookit/slog"
	"github.com/sudonym-btc/zap/service/config"
	gomail "gopkg.in/mail.v2"
)

func Connect(str string) (*gomail.Dialer, error) {
	parsedUrl, _ := url.Parse(str)
	password, _ := parsedUrl.User.Password()
	port, _ := strconv.Atoi(parsedUrl.Port())

	slog.Debug("Connecting to SMTP server", parsedUrl.Hostname(), port, parsedUrl.User.Username(), password)

	d := gomail.NewDialer(parsedUrl.Hostname(), port, parsedUrl.User.Username(), password)

	closer, err := d.Dial()
	if closer != nil {
		closer.Close()
	}
	if err != nil {
		slog.Warn("Failed connecting to SMTP server", err)
		return nil, err
	}
	slog.Debug("Connected to SMTP server")
	return d, nil
}

func SendMail(to string, subject string, body string) error {
	slog.Debug("Sending email", to, subject, body)

	conf, _ := config.LoadConfig()
	parsedUrl, _ := url.Parse(conf.Smtp)

	// Send email
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", parsedUrl.User.Username())

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	password, _ := parsedUrl.User.Password()
	port, _ := strconv.Atoi(parsedUrl.Port())

	d := gomail.NewDialer(parsedUrl.Hostname(), port, parsedUrl.User.Username(), password)

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		slog.Error("Failed sending email", err)

		return err
	}
	slog.Debug("Sent email")
	return nil

}
