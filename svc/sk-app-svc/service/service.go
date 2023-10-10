package service

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
}

type SkAppService struct {
}

func NewSkAppService() Service {
	return &SkAppService{}
}

func (svc *SkAppService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service
