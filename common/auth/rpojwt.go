package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct{
	ScreteKey string
}

type JwtClaims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

var (
	TokenExpire error = errors.New("token is expired")
	TokenRefresh error = errors.New("refresh default")
)

func (jw *JWT) CreateToken(clm JwtClaims) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256,clm)

	return token.SignedString(jw.ScreteKey)
}

func (jw *JWT) ParseToken(tokenString string) (*JwtClaims,error) {

	token, err := jwt.ParseWithClaims(tokenString,&JwtClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jw.ScreteKey,nil
	})
	if err != nil {
		switch e := err.(type) {
		case *jwt.ValidationError:
			switch e.Errors {
			case jwt.ValidationErrorExpired:
				return nil,TokenExpire
			default:
				break
			}
		}
		return nil,err
	}
	clms,ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return clms,nil
	}
	return nil,errors.New("unknown error")
}
// 在token 还没过期的时候刷新token， expireTime是unix过期时间戳
func (jw *JWT) RefreshToken(tokenString string,expireTime time.Time) (string,error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0,0)
	}

	token, err := jwt.ParseWithClaims(tokenString,&JwtClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jw.ScreteKey,nil
	})

	if err != nil {
		return "",err
	}

	if clm, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		clm.StandardClaims.ExpiresAt = expireTime.Unix()
		return jw.CreateToken(*clm)
	}
	return "",TokenRefresh
}

