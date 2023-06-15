package mail

import (
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/pkg/utils"
)

func SendEmail(emailTo dto.MailData, replyTo dto.MailData, data map[string]string, emailTemplate int64) error {
	hc := utils.NewHttpService("https://api.sendinblue.com/v3/")

	err := hc.SiBTempEmail(&utils.SendInBlueTempEmail{
		To: []utils.EmailContact{{Name: emailTo.Name, Email: emailTo.Email}},
		//Subject:    "New Order Demo",
		TemplateID: emailTemplate,
		Params:     data,
		ReplyTo:    utils.EmailContact{Name: replyTo.Name, Email: replyTo.Email},
		Sender:     utils.EmailContact{Name: replyTo.Name, Email: replyTo.Email},
	})

	return err
}
