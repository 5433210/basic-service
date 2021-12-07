package servicev1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apiv1courier "wailik.com/internal/courier/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

func SendEmail(s *service, toAddr string, code string) error {
	node := s.GetMicroService().Pick(constant.ServiceNameCourier)
	if node == nil {
		return errors.NewErrorC(errors.ErrCdCommon, nil)
	}

	log.Infof("node:%+v", node)

	client, err := apiv1courier.NewClient(node.Addr, apiv1courier.WithHTTPClient(http.DefaultClient))
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	// todo 方法链builder重构
	body := apiv1courier.SendEmailJSONRequestBody{}
	body.From = constant.EmailAddrNoReply
	body.To = make([]apiv1courier.EmailAddress, 1)
	body.To[0] = apiv1courier.EmailAddress(toAddr)
	template := apiv1courier.EmailTemplate{
		Id:   constant.EmailOtpTemplateId,
		Data: fmt.Sprintf(constant.EmailOtpTemplatePattern, code),
	}
	body.Content.Mode = constant.EmailModeTemplate
	body.Subject = constant.EmailOtpSubject

	payload, err := json.Marshal(&template)
	if err != nil {
		log.ErrorLog(err)

		return err
	}
	body.Content.Payload = payload

	log.Infof("payload:%+v", string(payload))

	resp, err := client.SendEmail(context.Background(), body)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	log.Infof("resp status:%+v", resp.Status)

	resp.Body.Close()

	return nil
}
