package core

import (
	"fmt"
	"github.com/ramzeng/user-microservice/pkg/config"
	"github.com/ramzeng/user-microservice/pkg/logger"
)

func initializeLogger() {
	var loggerConfig logger.Config

	if err := config.UnmarshalKey("logger", &loggerConfig); err != nil {
		panic(fmt.Sprintf("[core] | logger config unmarshal failed: %s", err))
	}

	if err := logger.Initialize(loggerConfig); err != nil {
		panic(fmt.Sprintf("[core] | logger initialize failed: %s", err))
	}
}
