package jwtAuth

import (
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	//Abstract Factory
	AuthFactory interface {
		SignToken() string
	}

	Claims struct {
		UserID uint `json:"user_id"`
		RoleID uint `josn:"role_id"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	//concrete
	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
)


func (a *authConcrete) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	ss, _ := token.SignedString(a.Secret)
	return ss
}

func now() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

// Note that: t is a second unit
func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	// 1 nanosec * 10^(9) = 1sec
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}

// int64 to jwtNumericDate adapter
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAccessToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &accessToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					// The "iss" (issuer) claim identifies the principal
					// that issued the JWT.  The processing of this claim is
					// generally application specific. (StringOrURI)
					Issuer: "book-store.com",

					// What is this token used for?
					Subject: "access-token",

					// Where you can use this token
					Audience: []string{"book-store.com"},

					ExpiresAt: jwtTimeDurationCal(expiredAt),

					// When this token is ready to use
					NotBefore: jwt.NewNumericDate(now()),

					// when you created this token
					IssuedAt: jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func NewRefreshToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "book-store.com",
					Subject:   "refresh-token",
					Audience:  []string{"book-store.com"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}
