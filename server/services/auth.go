package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/pkg/errors"
	"time"
)

type AuthService struct {
	ConfigRepository *config.Store
}

func NewAuthService(configRepository *config.Store) *AuthService {
	return &AuthService{
		configRepository,
	}
}

func (auth *AuthService) ChangePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password has to be at least 8 characters long")
	}

	hashedPassword := common.HashPassword(password)
	auth.ConfigRepository.Config.Set("password", hashedPassword)
	err := auth.ConfigRepository.Save()
	return err
}

func (auth *AuthService) Authorize(username string, password string) (string, error) {
	existingUsername := auth.ConfigRepository.Config.GetString("username")
	hashedPassword := auth.ConfigRepository.Config.GetString("password")
	secret := auth.ConfigRepository.Config.GetString("secret")
	if username != existingUsername || !common.IsCorrectPassword(password, hashedPassword) {
		return "", errors.New("incorrect username or password")
	}
	return getToken(username, secret)
}

func getToken(username string, secret string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["admin"] = true
	claims["created"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
