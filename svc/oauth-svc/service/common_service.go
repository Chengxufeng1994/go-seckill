package service

// CommonService Service Define a service interface
type CommonService interface {

	// HealthCheck check service health status
	HealthCheck() bool
}

type CommonServiceImpl struct {
}

func NewCommonService() CommonService {
	return &CommonServiceImpl{}
}

func (svc *CommonServiceImpl) HealthCheck() bool {
	return true
}
