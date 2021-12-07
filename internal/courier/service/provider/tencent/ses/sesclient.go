package ses

import (
	"encoding/json"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	apiv1 "wailik.com/internal/courier/api/v1"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/provider/tencent"
)

type sesClient struct {
	client *ses.Client
}

type SesClient interface {
	SendMail(email apiv1.Email) error
}

var _ SesClient = &sesClient{}

func New() (*sesClient, error) {
	credential := tencent.NewCredential()
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ses.tencentcloudapi.com"
	client, err := ses.NewClient(credential, "ap-hongkong", cpf)
	if err != nil {
		return nil, err
	}

	return &sesClient{client: client}, nil
}

func (c *sesClient) SendMail(email apiv1.Email) error {
	request := ses.NewSendEmailRequest()
	request.FromEmailAddress = (*string)(&email.From)
	for i := range email.To {
		request.Destination = append(request.Destination, (*string)(&email.To[i]))
	}
	request.Subject = &email.Subject
	if email.ReplyTo != nil && len(*(email.ReplyTo)) > 0 {
		var replyTo string
		for i := range *email.ReplyTo {
			replyTo += string((*email.ReplyTo)[i])
		}
	}

	var content map[string]interface{}
	err := json.Unmarshal([]byte(email.Content.Payload), &content)
	if err != nil {
		return err
	}
	log.Debugf("%+v", content)
	switch email.Content.Mode {
	case "template":
		var id uint64
		if id, err = strconv.ParseUint(content["id"].(string), 10, 64); err != nil {
			return err
		}
		log.Debugf("%+v", &id)
		data := content["data"].(string)
		log.Debugf("%+v", data)
		template := ses.Template{
			TemplateID:   &id,
			TemplateData: &data,
		}
		request.Template = &template
	case "text":
		text := content["text"].(string)
		request.Simple.Text = &text
	case "html":
		html := content["html"].(string)
		request.Simple.Html = &html
	}

	log.Debugf("request:%+v", request)
	_, err = c.client.SendEmail(request)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	log.Debugf("email sent")

	return nil
}
