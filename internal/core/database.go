package core

import (
	"fmt"
	"github.com/ramzeng/user-microservice/pkg/config"
	"github.com/ramzeng/user-microservice/pkg/database"
)

func initializeDatabase() {
	var databaseConfig database.Config

	if err := config.UnmarshalKey("mysql", &databaseConfig); err != nil {
		panic(fmt.Sprintf("[core] | database config unmarshal failed: %s", err))
	}

	if err := database.Initialize(databaseConfig); err != nil {
		panic(fmt.Sprintf("[core] | database initialize failed: %s", err))
	}
}
