package auth

import (
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (u *User) BeforeCreate(tx *sqlx.Tx) (err error) {
	u.ID = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	return nil
}

func (u *User) ComparePassword(hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(hash))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GenAuthToken() (*Auth, error) {
	accessTokenExp := time.Hour
	refreshTokenExp := time.Hour * 24 * 7
	accessToken, accessTokenUUID, err := makeToken(u, accessTokenExp)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshTokenUUID, err := makeToken(nil, refreshTokenExp)
	if err != nil {
		return nil, err
	}
	a := Auth{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int64(accessTokenExp.Minutes()),
		AccessTokenUUID:       accessTokenUUID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int64(accessTokenExp.Minutes()),
		RefreshTokenUUID:      refreshTokenUUID,
	}
	return &a, nil
}

type Auth struct {
	AccessToken           string
	AccessTokenExpiresIn  int64
	RefreshToken          string
	RefreshTokenExpiresIn int64
	AccessTokenUUID       string
	RefreshTokenUUID      string
}

type History struct {
	CreatedAt time.Time `json:"loginTime"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
}

func makeToken(u *User, exp time.Duration) (string, string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(exp).Unix(),
		Issuer:    "Pay9",
		IssuedAt:  time.Now().Unix(),
		Subject:   id.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", "", err
	}
	return tokenString, id.String(), nil
}
