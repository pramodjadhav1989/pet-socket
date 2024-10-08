package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/smartpet/websocket/models"
	"github.com/smartpet/websocket/utils/configs"
	log "github.com/smartpet/websocket/utils/logger"
)

var jwtSecretKey string

type Claims struct {
	jwt.StandardClaims
	OmneManagerID int16
	Token         string
	SourceID      string
}

func GenerateJWTAccessToken(userData models.TokenUserData) (string, error) {

	jwtSigninKey, err := configs.GetAppConfig("jwt_key", true)
	if err != nil {
		log.Error(context.Background()).Err(err).Msg(
			"error getting auth config")
		return "", err
	}
	//jwtSigninKey = configs.GetStringWithEnv(jwtSigninKey)
	secretkeyBytes := []byte(jwtSigninKey)

	tokenData := jwt.New(jwt.GetSigningMethod("HS256"))
	// prepare claims for token
	claims := models.JWTLoginToken{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "smartpet",
		},

		// add custom claims
		UserData: models.TokenUserData{
			CountryCode: userData.CountryCode,
			MobileNo:    userData.MobileNo,
			UserID:      userData.UserID,
			AppID:       userData.AppID,
			CreatedAt:   userData.CreatedAt,
		},
	}
	tokenData.Claims = claims

	// sign the generated key using secretKey
	token, err := tokenData.SignedString(secretkeyBytes)
	if err != nil {
		log.Error(context.Background()).Err(err).Msg(
			"error while signing the string")
		return "", err
	}
	return token, err
}

func ValidateJwtAndMatchClientIdCtx(ctx *http.Request, reqBodyclientId string) bool {

	err := AuthorizeSuperUser(ctx, reqBodyclientId) //validateJwtToken(bearerToken[7:], reqBodyclientId)

	if err != nil {
		return false
	}
	return true

}

func AuthorizeSuperUser(ctx *http.Request, partycode string) error {
	err := validateSupUserToken(ctx, partycode, "Authorization")
	if err == nil {
		return nil
	}

	err = validateSupUserToken(ctx, partycode, "AccessToken")
	if err == nil {
		return nil
	}

	err = validateSupUserToken(ctx, partycode, "Token")
	if err == nil {
		return nil
	}

	return err
}

func isSuperUserToken(clientcode string) (bool, error) {

	SuperUserIface, err := configs.GetAppConfig("superuserkey", true)
	if err != nil {
		return false, err
	}
	SuperUserKey := fmt.Sprintf("%v", SuperUserIface)

	log.Info(context.Background()).Msg("supper key " + SuperUserKey)

	return clientcode == SuperUserKey, nil

}

func validateSupUserToken(ctx *http.Request, partycode, header string) error {
	partycode = strings.ToUpper(partycode)
	auth := ctx.Header.Get(header)
	if auth == "" {
		return fmt.Errorf("empty header: %v", header)
	}
	authData, err := DecodeUserToken(auth)
	if err != nil {
		return err
	}
	clientcode := strings.ToUpper(strings.TrimSpace(authData.UserID))
	isSupe, err := isSuperUserToken(clientcode)
	if err != nil {
		return err
	}
	if clientcode != partycode && !isSupe {
		return fmt.Errorf("%v data not valid for partycode: %v", header, partycode)
	}
	return nil
}

func DecodeUserToken(tokenID string) (models.TokenUserData, error) {
	var (
		claim = &models.JWTLoginToken{}
		ok    bool
	)

	jwtSigninKey, err := configs.GetAppConfig("jwt_key", true)
	if err != nil {
		log.Error(context.Background()).Err(err).Msg(
			"error getting auth config")
		return models.TokenUserData{}, err
	}

	//jwtSigninKey = configs.GetStringWithEnv(jwtSigninKey)

	token, err := jwt.ParseWithClaims(tokenID, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSigninKey), nil
	})
	if err != nil {
		return models.TokenUserData{}, err
	}

	if claim, ok = token.Claims.(*models.JWTLoginToken); !ok || !token.Valid {
		return models.TokenUserData{}, errors.New("invalid token")
	}

	userData := claim.UserData

	if userData.UserID == "" && userData.MobileNo == "" {
		return models.TokenUserData{}, errors.New("user data not found in requested token")
	}
	return userData, nil
}
