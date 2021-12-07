package sms

import (
	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	apiv1 "wailik.com/internal/courier/api/v1"
	"wailik.com/internal/pkg/provider/tencent"
)

type smsClient struct {
	client *sms.Client
}

type SmsClient interface {
	SendSms(msg apiv1.Sms) error
}

var _ SmsClient = &smsClient{}

func New() (*smsClient, error) {
	credential := tencent.NewCredential()
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := sms.NewClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}

	return &smsClient{client: client}, nil
}

func (c *smsClient) SendSms(msg apiv1.Sms) error {
	request := sms.NewSendSmsRequest()

	for _, m := range msg.To {
		phoneNumber := fmt.Sprintf("%v+%v", (string)(m.Region), (string)(m.Number))
		request.PhoneNumberSet = append(request.PhoneNumberSet, &phoneNumber)
	}

	var content map[string]interface{}
	err := json.Unmarshal([]byte(msg.Content.Payload), &content)
	if err != nil {
		return err
	}
	switch msg.Content.Mode {
	case "template":
		id := content["id"].(string)
		params := content["params"].([]string)
		request.TemplateId = &id
		for i := range params {
			request.TemplateParamSet = append(request.TemplateParamSet, &params[i])
		}
	}

	_, err = c.client.SendSms(request)
	if err != nil {
		return err
	}

	return nil
}
