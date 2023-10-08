package token

import (
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type OAuth2TokenCustomClaims struct {
	UserDetails   model.UserDetails
	ClientDetails model.ClientDetails
	RefreshToken  model.OAuth2Token
	jwt.StandardClaims
}

type JwtTokenEnhancer struct {
	secretKey []byte
}

func NewJwtTokenEnhancer(secretKey []byte) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: secretKey,
	}
}

func (enhancer *JwtTokenEnhancer) Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return enhancer.sign(oauth2Token, oauth2Details)
}

func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return enhancer.secretKey, nil
	}
	jwtToken, err := jwt.ParseWithClaims(tokenValue, &OAuth2TokenCustomClaims{}, keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims := jwtToken.Claims.(*OAuth2TokenCustomClaims)

	expiresTime := time.Unix(claims.ExpiresAt, 0)

	return &model.OAuth2Token{
			RefreshToken: &claims.RefreshToken,
			TokenValue:   tokenValue,
			ExpiresTime:  &expiresTime,
		}, &model.OAuth2Details{
			User:   &claims.UserDetails,
			Client: &claims.ClientDetails,
		}, nil
}

func (enhancer *JwtTokenEnhancer) sign(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	expireTime := oauth2Token.ExpiresTime
	clientDetails := *oauth2Details.Client
	userDetails := *oauth2Details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	claims := OAuth2TokenCustomClaims{
		UserDetails:   userDetails,
		ClientDetails: clientDetails,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "System",
		},
	}

	if oauth2Token.RefreshToken != nil {
		claims.RefreshToken = *oauth2Token.RefreshToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := token.SignedString(enhancer.secretKey)
	if err != nil {
		return nil, err
	}

	oauth2Token.TokenValue = tokenValue
	oauth2Token.TokenType = "jwt"
	return oauth2Token, nil
}
