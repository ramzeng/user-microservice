package core

import (
	"fmt"
	"github.com/ramzeng/user-microservice/pkg/config"
)

func ReadConfig(path string) {
	if err := config.Initialize(path); err != nil {
		panic(fmt.Sprintf("[core] | config file read failed: %s", err))
	}
}
