package email

import (
	"context"
	"errors"

	"github.com/mailgun/mailgun-go/v3"
)

type Attachment struct {
	FileName string
	Content  []byte
}
type SendMsg struct {
	Body    string       `json:"body"`
	To      []string     `json:"to"`
	Cc      []string     `json:"cc"`  // 抄送
	Bcc     []string     `json:"bcc"` // 秘密抄送
	Subject string       `json:"subject"`
	Attach  []Attachment `json:"attach"` // 附件
}

type MailGunConfig struct {
	MG *mailgun.MailgunImpl
}

func New(domain, apiKey string) *MailGunConfig {
	return &MailGunConfig{
		MG: mailgun.NewMailgun(domain, apiKey),
	}
}
func (m *MailGunConfig) SendMailGun(msg *SendMsg, sender string) error {
	if len(msg.To) == 0 {
		return errors.New("email receiver is empty")
	}
	message := m.MG.NewMessage(sender, msg.Subject, "", msg.To...)
	message.SetHtml(msg.Body)

	for _, cc := range msg.Cc {
		message.AddCC(cc)
	}

	for _, bcc := range msg.Bcc {
		message.AddBCC(bcc)
	}

	for _, attach := range msg.Attach {
		message.AddBufferAttachment(attach.FileName, attach.Content)
	}

	_, _, err := m.MG.Send(context.Background(), message)
	return err
}
