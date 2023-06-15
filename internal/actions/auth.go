package actions

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

type SignedClaim struct {
	User *dto.UserJwtEntity
	Uid  string
	jwt.StandardClaims
}

type TokenDetail struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func PasswordBcrypt(password string) (string, error) {
	bp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bp), nil
}

func VerifyPassword(pwd, hashedPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

func IssueToken(user *dto.UserJwtEntity) (*TokenDetail, error) {
	accessTokenAlive := configs.AppConfig.AccessTokenAlive
	aToken, _ := strconv.ParseInt(accessTokenAlive, 10, 64)
	refreshTokenAlive := configs.AppConfig.RefreshTokenAlive
	rToken, _ := strconv.ParseInt(refreshTokenAlive, 10, 64)

	td := &TokenDetail{}
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(aToken)).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Minute * time.Duration(rToken)).Unix()
	td.RefreshUuid = uuid.New().String()
	user.Asid = td.AccessUuid
	user.Rsid = td.RefreshUuid

	secret, err := signedPrivateSecret("id_rsa")
	if err != nil {
		return nil, err
	}

	aClaims := &SignedClaim{
		User: user,
		Uid:  td.AccessUuid,
		StandardClaims: jwt.StandardClaims{
			Id:        td.AccessUuid,
			Issuer:    configs.AppConfig.AppKey,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: td.AtExpires,
		},
	}

	at, err := jwt.NewWithClaims(jwt.SigningMethodRS256, aClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}
	td.AccessToken = at

	//Generating refresh token
	rSecret, err := signedPrivateSecret("id_rsa")
	if err != nil {
		return nil, err
	}
	rUser := user
	rUser.Rsid = td.AccessUuid
	rUser.Rsid = td.RefreshUuid

	rClaims := &SignedClaim{
		User: rUser,
		Uid:  td.RefreshUuid,
		StandardClaims: jwt.StandardClaims{
			Id:        td.RefreshUuid,
			Issuer:    configs.AppConfig.AppName,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: td.RtExpires,
		},
	}
	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rClaims).SignedString(rSecret)
	if err != nil {
		return nil, err
	}
	td.RefreshToken = rt
	return td, nil
}

func ValidateToken(jwtToken string) (*jwt.Token, error) {
	key, err := signedPublicSecret("id_rsa.pub")
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	return jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
}

func ValidateRefreshToken(refreshToken string) (*jwt.Token, error) {
	key, err := signedPublicSecret("id_rsa.pub")
	if err != nil {
		return nil, err
	}
	return jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
}

func signedPublicSecret(secretFileName string) (*rsa.PublicKey, error) {
	secretFileData, err := readRsaFile(secretFileName)
	if err != nil {
		return nil, err
	}
	secret, err := jwt.ParseRSAPublicKeyFromPEM([]byte(secretFileData))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func signedPrivateSecret(secretFileName string) (*rsa.PrivateKey, error) {
	secretFileData, err := readRsaFile(secretFileName)
	if err != nil {
		return nil, err
	}
	secret, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secretFileData))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func readRsaFile(filename string) (string, error) {
	data, err := os.ReadFile("./certs/" + filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
