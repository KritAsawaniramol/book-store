package jwtAuth

import (
	"errors"
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

	// custom claims
	AuthMapClaims struct {
		*Claims
		// standard claims
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
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ReloadToken(secret string, expiredAt int64, claims *Claims) string {
	obj := &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "myAnimeList.com",
					Subject:   "refresh-token",
					Audience:  []string{"myAnimeList.com"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
	return obj.SignToken()
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthMapClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// check this token sighn with expected method ( SignToken() => jwt.SigningMethodHS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("error: unexpected signing method")
			}
			return []byte(secret), nil
		})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("error: token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("error: token is expired")
		} else {
			return nil, errors.New("error: token is invalid")
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error: claims type is invalid")
	}
}
