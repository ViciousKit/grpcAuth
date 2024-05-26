package sso

import (
	"sso/internal/config"

	sso_v1 "github.com/ViciousKit/proto/generated/go/sso"
)

func main() {
	sso_v1.RegisterAuthServer()
	c := config.MustLoad()
	//TODO: initialize cfg

	//TODO: initialize logger

	//TODO: initialize app

	//TODO: run grpc server
}
