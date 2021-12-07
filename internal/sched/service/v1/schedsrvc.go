package servicev1

const (
	ModeSimple = "simple"
)

type schedSrvc struct {
	service *service
}

func newSchedSrvc(s *service) *schedSrvc {
	return &schedSrvc{
		service: s,
	}
}
