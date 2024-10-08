package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Level int8

const (
	// InfoLevel defines info log level.
	InfoLevel Level = iota
	// WarnLevel defines warn log level.
	WarnLevel
)

type LogMessage struct {
	LogLevel Level
	Message  string
}

type JWTTokenClaims struct {
	UserData      TokenUserData `json:"userData,omitempty"`
	OmneManagerID int16         `json:"omnemanagerid,omitempty"`
	Token         string        `json:"token,omitempty"`
	SourceID      string        `json:"sourceid,omitempty"`
	TokenClaims
	jwt.RegisteredClaims
}

type TokenUserData struct {
	CountryCode string    `json:"country_code,omitempty"`
	MobileNo    string    `json:"mob_no,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Source      string    `json:"source,omitempty"`
	AppID       string    `json:"app_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	DataCenter  string    `json:"dataCenter,omitempty"`
}

type SmartPetClaims struct {
	UserType   string           `json:"user_type,omitempty"` // accepted value application, client, admin
	MobileNo   string           `json:"mobile_no,omitempty"`
	TokenType  string           `json:"token_type,omitempty"` // accepted value trade_access_token, trade_refresh_token, non_trade_access_token, non_trade_refresh_token
	DataCenter string           `json:"data_center,omitempty"`
	GMId       string           `json:"gm_id,omitempty"`
	Source     string           `json:"source,omitempty"`
	DeviceId   string           `json:"device_id,omitempty"`
	Issuer     string           `json:"issuer"`
	Subject    string           `json:"subject,omitempty"`
	Audience   jwt.ClaimStrings `json:"audience,omitempty"`
	Scope      jwt.ClaimStrings `json:"scope,omitempty"` // will be used for authorization roles
	KeyId      string           `json:"key_id"`
}

type TokenClaims struct {
	UserType   string           `json:"user_type,omitempty"` // accepted value application, client, admin
	MobileNo   string           `json:"mobile_no,omitempty"`
	TokenType  string           `json:"token_type,omitempty"` // accepted value trade_access_token, trade_refresh_token, non_trade_access_token, non_trade_refresh_token
	Scope      jwt.ClaimStrings `json:"scope,omitempty"`      // will be used for authorization roles
	DataCenter string           `json:"data_center,omitempty"`
	GMId       string           `json:"gm_id,omitempty"`
	Source     string           `json:"source,omitempty"`
	DeviceId   string           `json:"device_id,omitempty"`
	KeyId      string           `json:"kid,omitempty"`
}

type JWTLoginToken struct {
	UserData TokenUserData `json:"userData"`
	jwt.StandardClaims
}

type Actor struct {
	Subject string `json:"sub,omitempty"`
}

func (c JWTTokenClaims) Valid() error {
	return c.RegisteredClaims.Valid()
}

type UnAuthorizeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
