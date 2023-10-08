package entity

import (
	"gorm.io/gorm"
	"time"
)

type ClientDetails struct {
	ID        uint `gorm:"auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// client identity
	ClientId string `gorm:"size:255;primary_key;not null"`
	// client secret
	ClientSecret string `gorm:"size:255;not null"`
	// access token validity seconds
	AccessTokenValiditySeconds int `gorm:"size:10;not null"`
	// refresh token validity seconds
	RefreshTokenValiditySeconds int `gorm:"size:10;not null"`
	// redirect uri
	RegisteredRedirectUri string `gorm:"size:128;not null"`
	// grant types
	AuthorizedGrantTypes string `gorm:"size:128;not null"`
}
