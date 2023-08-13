package auth

type Config struct {
	Secret                 string
	AccessTokenTimeToLive  int
	RefreshTokenTimeToLive int
}
