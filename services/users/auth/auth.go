package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mirzakhany/pm/pkg/config"
	"github.com/mirzakhany/pm/pkg/session"
	users "github.com/mirzakhany/pm/services/users/proto"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret        = config.RegisterString("auth.jwtSecret", "this-is-for-test-dont-use-in-production")
	accessTokenLife  = config.RegisterInt("auth.accessTokenLife", 15)
	refreshTokenLife = config.RegisterInt("auth.refreshTokenLife", 170)
)

// SaveTokens save user tokens after login
func SaveTokens(user *users.User, response *users.LoginResponse) error {
	err := session.Set(response.AccessToken, user, time.Minute*time.Duration(accessTokenLife.Int()))
	if err != nil {
		return err
	}
	return session.Set(response.RefreshToken, user, time.Hour*time.Duration(refreshTokenLife.Int()))
}

// CreateToken will create access and refresh token
func CreateToken(user *users.User) (*users.LoginResponse, error) {

	accessExpiresAt := time.Now().Add(time.Minute * time.Duration(accessTokenLife.Int())).Unix()
	refreshExpiresAt := time.Now().Add(time.Hour * time.Duration(refreshTokenLife.Int())).Unix()
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.Id,
		"access_uuid": user.Uuid,
		"exp":         accessExpiresAt,
	})

	accessToken, err := aToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      user.Id,
		"refresh_uuid": user.Uuid,
		"exp":          refreshExpiresAt,
	})

	refreshToken, err := rToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	return &users.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// HashPassword return hashed password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash will check hashed password against password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
