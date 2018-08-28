package token

import (
	"fmt"
	log "questionair_backend/util/logger"
	"time"

	"questionair_backend/conf"
	"questionair_backend/defines"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func CreateToken(info map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range info {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(time.Duration(conf.Conf.Token.TokenExpire) * time.Second).Unix()
	t, err := token.SignedString([]byte(conf.Conf.Token.Salt))
	if err != nil {
		return "", err
	}
	return t, nil
}

func ParseToken(tokenString string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.Token.Salt), nil
	})
	if err != nil {
		return res, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for k, v := range claims {
			res[k] = v
		}
		return res, nil
	}
	return res, fmt.Errorf("invalid token")
}

func GetTokenId(e echo.Context) (int64, error) {
	user := e.Get("usr").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	if id == 0 {
		return 0, defines.ComAuthFailed
	}
	return int64(id), nil
}

func GetTokenRole(e echo.Context) (string, error) {
	user := e.Get("usr").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role == "" {
		return "", defines.ComAuthFailed
	}
	return role, nil
}

// verifyUser vertify admin; check role and id
func VerifyUser(e echo.Context) (bool, error) {
	if role, err := GetTokenRole(e); err != nil || role != "user" {
		log.Logger().Errorf("VerifyUser: check user failed")
		return false, defines.ComAuthFailed
	}
	if id, err := GetTokenId(e); err != nil || id == 0 {
		log.Logger().Errorf("VerifyUser: check user failed")
		return false, defines.ComAuthFailed
	}
	return true, nil
}
