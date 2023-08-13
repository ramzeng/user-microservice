package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type Token struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiredAt  int64  `json:"access_token_expired_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredAt int64  `json:"refresh_token_expired_at"`
}

func CreateToken(subject string, withRefreshToken bool) (Token, error) {
	token := Token{}

	now := time.Now()
	uuidString := uuid.NewString()

	accessTokenExpiredAt := now.Add(time.Second * time.Duration(config.AccessTokenTimeToLive)).Unix()

	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: accessTokenExpiredAt,
		Id:        uuidString,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Subject:   subject,
	})

	if accessToken, err := accessTokenClaims.SignedString([]byte(config.Secret)); err != nil {
		return token, err
	} else {
		token.AccessToken = accessToken
		token.AccessTokenExpiredAt = accessTokenExpiredAt
	}

	if !withRefreshToken {
		return token, nil
	}

	refreshTokenExpiredAt := now.Add(time.Second * time.Duration(config.RefreshTokenTimeToLive)).Unix()

	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: refreshTokenExpiredAt,
		Id:        uuidString,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Subject:   subject,
	})

	if refreshToken, err := refreshTokenClaims.SignedString([]byte(config.Secret)); err != nil {
		return token, err
	} else {
		token.RefreshToken = refreshToken
		token.RefreshTokenExpiredAt = refreshTokenExpiredAt
	}

	return token, nil
}

func RefreshAccessToken(accessToken string) (Token, error) {
	accessTokenClaims, err := ParseTokenClaims(accessToken)

	if err != nil {
		return Token{}, err
	}

	if accessTokenClaims.ExpiresAt > time.Now().Unix() {
		return Token{
			AccessToken:          accessToken,
			AccessTokenExpiredAt: accessTokenClaims.ExpiresAt,
		}, nil
	}

	now := time.Now()

	refreshedAccessTokenExpiredAt := now.Add(time.Second * time.Duration(config.AccessTokenTimeToLive)).Unix()

	refreshedAccessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: refreshedAccessTokenExpiredAt,
		Id:        accessTokenClaims.Id,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Subject:   accessTokenClaims.Subject,
	})

	refreshedAccessToken, err := refreshedAccessTokenClaims.SignedString([]byte(config.Secret))

	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken:          refreshedAccessToken,
		AccessTokenExpiredAt: refreshedAccessTokenExpiredAt,
	}, nil
}

func ParseToken(inputToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(inputToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Secret), nil
	})
}

func ParseTokenClaims(inputToken string) (*jwt.StandardClaims, error) {
	token, _ := ParseToken(inputToken)

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func ParseTokenSubject(inputToken string) (string, error) {
	token, err := ParseToken(inputToken)

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims.Subject, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func VerifyToken(inputToken string) bool {
	token, err := ParseToken(inputToken)

	if err != nil {
		return false
	}

	if _, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return true
	} else {
		return false
	}
}

func MatchAccessTokenAndRefreshToken(accessToken string, refreshToken string) bool {
	accessTokenClaims, err := ParseTokenClaims(accessToken)

	if err != nil {
		return false
	}

	refreshTokenClaims, err := ParseTokenClaims(refreshToken)

	if err != nil {
		return false
	}

	if accessTokenClaims.Id != refreshTokenClaims.Id {
		return false
	}

	if accessTokenClaims.Subject != refreshTokenClaims.Subject {
		return false
	}

	return true
}
