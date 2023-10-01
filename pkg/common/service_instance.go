package common

type ServiceInstance struct {
	Host      string
	Port      int
	Weight    int
	CurWeight int
	GrpcPort  int
}
