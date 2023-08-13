package core

import (
	"fmt"
	"github.com/ramzeng/user-microservice/pkg/auth"
	"github.com/ramzeng/user-microservice/pkg/config"
)

func initializeToken() {
	var authConfig auth.Config

	if err := config.UnmarshalKey("auth", &authConfig); err != nil {
		panic(fmt.Sprintf("[core] | parse auth config failed: %s", err))
	}

	auth.Initialize(authConfig)
}
