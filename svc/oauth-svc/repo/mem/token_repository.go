package mem

import (
	"errors"
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/token"
	"sync"
)

var ErrNotSupportOperation = errors.New("no support operation")

type TokenStore interface {
	StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
	ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error)
	GetAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	RemoveAccessToken(tokenValue string)
	StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	RemoveRefreshToken(oauth2Token string)
	ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error)
	ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error)
}

type InMemTokenRepository struct {
	enhancer       token.TokenEnhancer
	oauth2TokenMap sync.Map
	mu             sync.RWMutex
}

func NewInMemTokenRepository(enhancer token.TokenEnhancer) TokenStore {
	return &InMemTokenRepository{
		enhancer: enhancer,
	}
}

func (r *InMemTokenRepository) StoreAccessToken(oauthToken *model.OAuth2Token, oauthDetails *model.OAuth2Details) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("access_token:%s:%d", oauthDetails.Client.ClientId, oauthDetails.User.UserId)
	r.oauth2TokenMap.Store(key, oauthToken)
	//time.AfterFunc(ttl, func() {
	//	r.oauth2TokenMap.Delete(key)
	//})
}

func (r *InMemTokenRepository) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	oAuth2Token, _, err := r.enhancer.Extract(tokenValue)
	if err != nil {
		return nil, err
	}
	return oAuth2Token, nil
}

func (r *InMemTokenRepository) ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, oAuthDetails, err := r.enhancer.Extract(tokenValue)
	if err != nil {
		return nil, err
	}
	return oAuthDetails, nil
}

func (r *InMemTokenRepository) GetAccessToken(oauthDetails *model.OAuth2Details) (*model.OAuth2Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("access_token:%s:%d", oauthDetails.Client.ClientId, oauthDetails.User.UserId)
	value, ok := r.oauth2TokenMap.Load(key)
	if !ok {
		return nil, ErrNotSupportOperation
	}
	return value.(*model.OAuth2Token), nil
}

func (r *InMemTokenRepository) RemoveAccessToken(tokenValue string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, oauthDetails, _ := r.enhancer.Extract(tokenValue)
	key := fmt.Sprintf("access_token:%s:%d", oauthDetails.Client.ClientId, oauthDetails.User.UserId)
	r.oauth2TokenMap.Delete(key)
}

func (r *InMemTokenRepository) StoreRefreshToken(oauthToken *model.OAuth2Token, oauthDetails *model.OAuth2Details) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("refresh_token:%s:%d", oauthDetails.Client.ClientId, oauthDetails.User.UserId)
	r.oauth2TokenMap.Store(key, oauthToken)
}

func (r *InMemTokenRepository) RemoveRefreshToken(tokenValue string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, oauthDetails, _ := r.enhancer.Extract(tokenValue)
	key := fmt.Sprintf("refresh_token:%s:%d", oauthDetails.Client.ClientId, oauthDetails.User.UserId)
	r.oauth2TokenMap.Delete(key)
}

func (r *InMemTokenRepository) ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	oAuth2Token, _, err := r.enhancer.Extract(tokenValue)
	if err != nil {
		return nil, err
	}
	return oAuth2Token, nil
}

func (r *InMemTokenRepository) ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error) {
	_, oauth2Details, err := r.enhancer.Extract(tokenValue)
	return oauth2Details, err
}
