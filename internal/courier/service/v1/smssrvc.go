package servicev1

import apiv1 "wailik.com/internal/courier/api/v1"

type smsSrvc struct {
	service *service
}

func newSmsSrvc(s *service) *smsSrvc {
	return &smsSrvc{service: s}
}

func (s *smsSrvc) Send(sms apiv1.Sms) error {
	return nil
}
