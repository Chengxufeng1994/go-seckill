package main

import "github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/setup"

func main() {
	setup.InitializeZk()
	setup.InitializeRedis()
	setup.InitializeService()
}
