package service

type CommonService interface {
	HealthCheck() bool
}

type CommonServiceImpl struct {
}

func NewCommonService() CommonService {
	return &CommonServiceImpl{}
}

func (s *CommonServiceImpl) HealthCheck() bool {
	return true
}

type CommonServiceMiddleware func(CommonService) CommonService
