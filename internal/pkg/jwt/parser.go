package jwt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"

	"github.com/cage1016/ms-sample/internal/pkg/errors"
)

type contextKey string

const (
	XJWTClaimsContextKey contextKey = "XJWTClaims"
)

var (
	ErrXJWTContextMissing = errors.New("xjp up for parsing was not passed through the context")

	ErrXJWTBase64DecodeFail = errors.New("xjp base64 decode error")

	ErrClaimsInvalid = errors.New("Claims is invalid")
)

func XJWTParser() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			payload, ok := ctx.Value(XJwtPlayload).(string)
			if !ok {
				return nil, ErrXJWTContextMissing
			}

			xp, err := base64.RawURLEncoding.DecodeString(payload)
			if err != nil {
				return nil, ErrXJWTBase64DecodeFail
			}

			var cl stdjwt.MapClaims
			if err = json.NewDecoder(bytes.NewBuffer(xp)).Decode(&cl); err != nil {
				return nil, ErrClaimsInvalid
			}

			ctx = context.WithValue(ctx, XJWTClaimsContextKey, cl)
			return next(ctx, request)
		}
	}
}
