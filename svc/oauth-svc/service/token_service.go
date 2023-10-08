package service

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/repo/mem"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/token"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrNotSupportGrantType               = errors.New("grant type is not supported")
	ErrNotSupportOperation               = errors.New("no support operation")
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username, password")
	ErrInvalidTokenRequest               = errors.New("invalid token")
	ErrExpiredToken                      = errors.New("token is expired")
)

type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error)
}

type ComposeTokenGranter struct {
	TokenGrantDict map[string]TokenGranter
}

func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{
		TokenGrantDict: tokenGrantDict,
	}
}

func (tg *ComposeTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	dispatchGranter, ok := tg.TokenGrantDict[grantType]
	if !ok {
		return nil, ErrNotSupportGrantType
	}

	return dispatchGranter.Grant(ctx, grantType, client, reader)
}

type UsernamePasswordTokenGranter struct {
	supportGrantType   string
	tokenService       TokenService
	userDetailsService UserDetailsService
}

func NewUsernamePasswordTokenGranter(grantType string, userDetailsService UserDetailsService, tokenService TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGrantType:   grantType,
		tokenService:       tokenService,
		userDetailsService: userDetailsService,
	}
}

func (tg *UsernamePasswordTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tg.supportGrantType {
		return nil, ErrNotSupportGrantType
	}

	username := reader.FormValue("username")
	password := reader.FormValue("password")
	if username == "" || password == "" {
		return nil, ErrInvalidUsernameAndPasswordRequest
	}

	user, err := tg.userDetailsService.GetUserDetailByUsername(ctx, username, password)
	if err != nil {
		return nil, ErrInvalidUsernameAndPasswordRequest
	}

	return tg.tokenService.CreateAccessToken(
		&model.OAuth2Details{
			Client: client,
			User:   user,
		})
}

type RefreshTokenGranter struct {
	supportGrantType string
	tokenService     TokenService
}

func NewRefreshGranter(grantType string, userDetailsService UserDetailsService, tokenService TokenService) TokenGranter {
	return &RefreshTokenGranter{
		supportGrantType: grantType,
		tokenService:     tokenService,
	}
}

func (tg *RefreshTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tg.supportGrantType {
		return nil, ErrNotSupportGrantType
	}

	refreshTokenValue := reader.URL.Query().Get("refresh_token")

	if refreshTokenValue == "" {
		return nil, ErrInvalidTokenRequest
	}

	return tg.tokenService.RefreshAccessToken(refreshTokenValue)

}

type TokenService interface {
	// 根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	// 根据用户信息和客户端信息生成访问令牌
	CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	// 根据用户信息和客户端信息获取已生成访问令牌
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据访问令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
}

type DefaultTokenService struct {
	tokenEnhancer token.TokenEnhancer
	repo          mem.TokenStore
}

func NewTokenService(tokenEnhancer token.TokenEnhancer, repository mem.TokenStore) TokenService {
	return &DefaultTokenService{
		tokenEnhancer: tokenEnhancer,
		repo:          repository,
	}
}

func (svc *DefaultTokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error) {
	accessToken, err := svc.ReadAccessToken(tokenValue)
	if err != nil {
		return nil, err
	}

	if accessToken.IsExpired() {
		return nil, ErrExpiredToken
	}

	return svc.repo.ReadOAuth2Details(tokenValue)
}

func (svc *DefaultTokenService) _createAccessToken(refreshToken *model.OAuth2Token, details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validaitySeconds := details.Client.RefreshTokenValiditySeconds
	duration, _ := time.ParseDuration(strconv.Itoa(validaitySeconds) + "s")
	expiresTime := time.Now().Add(duration)
	accessToken := &model.OAuth2Token{
		RefreshToken: refreshToken,
		TokenValue:   uuid.NewV4().String(),
		ExpiresTime:  &expiresTime,
	}

	if svc.tokenEnhancer != nil {
		return svc.tokenEnhancer.Enhance(accessToken, details)
	}

	return accessToken, nil
}

func (svc *DefaultTokenService) CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	var refreshToken *model.OAuth2Token
	existedToken, err := svc.GetAccessToken(oauth2Details)
	if existedToken != nil {
		if existedToken.IsExpired() {
			svc.repo.RemoveAccessToken(existedToken.TokenValue)
			if existedToken.RefreshToken != nil {
				refreshToken = existedToken.RefreshToken
				svc.repo.RemoveRefreshToken(refreshToken.TokenValue)
			}
		} else {
			svc.repo.StoreAccessToken(existedToken, oauth2Details)
			return existedToken, nil
		}
	}

	if refreshToken == nil || refreshToken.IsExpired() {
		refreshToken, err = svc._createRefreshToken(oauth2Details)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := svc._createAccessToken(refreshToken, oauth2Details)
	if err != nil {
		return accessToken, err
	}

	svc.repo.StoreAccessToken(accessToken, oauth2Details)
	svc.repo.StoreRefreshToken(refreshToken, oauth2Details)

	return accessToken, err
}

func (svc *DefaultTokenService) _createRefreshToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validaitySeconds := details.Client.RefreshTokenValiditySeconds
	duration, _ := time.ParseDuration(strconv.Itoa(validaitySeconds) + "s")
	expiresTime := time.Now().Add(duration)
	refreshToken := &model.OAuth2Token{
		TokenValue:  uuid.NewV4().String(),
		ExpiresTime: &expiresTime,
	}

	if svc.tokenEnhancer != nil {
		return svc.tokenEnhancer.Enhance(refreshToken, details)
	}

	return refreshToken, nil
}

func (svc *DefaultTokenService) RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error) {
	refreshToken, err := svc.repo.ReadRefreshToken(refreshTokenValue)
	if err != nil {
		return nil, err
	}
	if refreshToken.IsExpired() {
		return nil, ErrExpiredToken
	}

	oauth2Details, err := svc.repo.ReadOAuth2DetailsForRefreshToken(refreshTokenValue)
	if err != nil {
		return nil, err
	}

	oauth2Token, err := svc.repo.GetAccessToken(oauth2Details)
	if err != nil {
		return nil, err
	}
	svc.repo.RemoveAccessToken(oauth2Token.TokenValue)
	svc.repo.RemoveRefreshToken(refreshTokenValue)
	refreshToken, err = svc._createRefreshToken(oauth2Details)
	if err != nil {
		return nil, err
	}
	accessToken, err := svc._createAccessToken(refreshToken, oauth2Details)
	if err != nil {
		return nil, err
	}

	svc.repo.StoreAccessToken(accessToken, oauth2Details)
	svc.repo.StoreAccessToken(refreshToken, oauth2Details)

	return accessToken, nil
}

func (svc *DefaultTokenService) GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return svc.repo.GetAccessToken(details)
}

func (svc *DefaultTokenService) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	return svc.repo.ReadAccessToken(tokenValue)
}
