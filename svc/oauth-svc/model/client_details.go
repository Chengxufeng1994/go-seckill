package model

import (
	"encoding/json"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/entity"
)

type ClientDetails struct {
	ID uint
	// client identity
	ClientId string
	// client secret
	ClientSecret string
	// access token validity seconds
	AccessTokenValiditySeconds int
	// refresh token validity seconds
	RefreshTokenValiditySeconds int
	// redirect uri
	RegisteredRedirectUri string
	// grant types
	AuthorizedGrantTypes []string
}

func ClientDetailsModel2Entity(model *ClientDetails) (*entity.ClientDetails, error) {
	authorizedGrantTypes, err := json.Marshal(&model.AuthorizedGrantTypes)
	if err != nil {
		return nil, err
	}

	return &entity.ClientDetails{
		ID:                          model.ID,
		ClientId:                    model.ClientId,
		ClientSecret:                model.ClientSecret,
		AccessTokenValiditySeconds:  model.AccessTokenValiditySeconds,
		RefreshTokenValiditySeconds: model.RefreshTokenValiditySeconds,
		RegisteredRedirectUri:       model.RegisteredRedirectUri,
		AuthorizedGrantTypes:        string(authorizedGrantTypes),
	}, nil
}

func ClientDetailsEntity2Model(entity *entity.ClientDetails) (*ClientDetails, error) {
	var authorizedGrantTypes []string
	err := json.Unmarshal([]byte(entity.AuthorizedGrantTypes), &authorizedGrantTypes)
	if err != nil {
		return nil, err
	}

	return &ClientDetails{
		ClientId:                    entity.ClientId,
		ClientSecret:                entity.ClientSecret,
		AccessTokenValiditySeconds:  entity.AccessTokenValiditySeconds,
		RefreshTokenValiditySeconds: entity.RefreshTokenValiditySeconds,
		RegisteredRedirectUri:       entity.RegisteredRedirectUri,
		AuthorizedGrantTypes:        authorizedGrantTypes,
	}, nil
}
