package database

type Config struct {
	User                         string
	Password                     string
	Host                         string
	Port                         string
	Database                     string
	MaxIdleConnections           int
	MaxOpenConnections           int
	ConnectionMaxIdleSeconds     int
	ConnectionMaxLifetimeSeconds int
}
