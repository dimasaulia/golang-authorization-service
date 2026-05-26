package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type Mailer interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type SMTPMailer struct {
	cfg config.MailConfig
	log *logger.LayerLogger
}

func New(cfg config.Config, appLogger *logger.Logger) Mailer {
	return &SMTPMailer{
		cfg: cfg.Mail,
		log: appLogger.Layer("platform.mailer"),
	}
}

func (m *SMTPMailer) Send(ctx context.Context, to string, subject string, body string) error {
	if !m.cfg.Enabled {
		m.log.Info(ctx, "send.skipped", "to", to, "subject", subject)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	message := strings.Join([]string{
		"From: " + m.cfg.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, []byte(message)); err != nil {
		m.log.Error(ctx, "send.failed", err, "to", to, "subject", subject)
		return err
	}

	m.log.Info(ctx, "send.success", "to", to, "subject", subject)
	return nil
}
