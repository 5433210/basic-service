package servicev1

import (
	apiv1 "wailik.com/internal/courier/api/v1"
	"wailik.com/internal/courier/service/provider/tencent/ses"
)

type emailSrvc struct {
	service   *service
	sesclient ses.SesClient
}

func newEmailSrvc(s *service) *emailSrvc {
	sesClient, err := ses.New()
	if err != nil {
		return nil
	}

	return &emailSrvc{
		service:   s,
		sesclient: sesClient,
	}
}

func (s *emailSrvc) Send(email apiv1.Email) (*[]apiv1.SendStatus, error) {
	// no send status in ses
	status := make([]apiv1.SendStatus, 0)
	err := s.sesclient.SendMail(email)
	if err != nil {
		return nil, err
	}

	return &status, nil
}
