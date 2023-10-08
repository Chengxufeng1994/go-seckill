package token

import "github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"

type TokenEnhancer interface {
	// Enhance to create token with info
	Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// Extract get info from token
	Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error)
}
