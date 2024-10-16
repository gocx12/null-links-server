package internal

import "github.com/golang-jwt/jwt"

// @secretKey: JWT secret key
// @iat: time stamp
// @seconds: expire time(second)
// @payload: data payload
func GenJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func ParseJwtToken(secretKey, token string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return result, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return result, err
	}

	for k, v := range claims {
		result[k] = v
	}
	return result, nil
}
