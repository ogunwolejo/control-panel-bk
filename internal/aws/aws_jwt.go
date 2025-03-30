package aws

import (
	"control-panel-bk/config"
	"log"
	"time"
)

type CognitoToken struct {
	IdToken      string
	AccessToken  string // An hour live span
	RefreshToken string // A min of 30 days to a max of 10 years life span
	TokenType *string
	ExpiresIn time.Time
}

func (c *CognitoToken) DecodeIdToken() {}

func (c *CognitoToken) VerifyIdToken() {}

func (c *CognitoToken) RefreshingSessionToken(clientId string) error {
	tokens, err := GetRefreshToken(config.AwsConfig, clientId, c.RefreshToken)
	if err != nil {
		return err
	}

	c.IdToken = *tokens.AuthenticationResult.IdToken
	c.AccessToken = *tokens.AuthenticationResult.AccessToken
	c.RefreshToken = *tokens.AuthenticationResult.RefreshToken
	c.TokenType = tokens.AuthenticationResult.TokenType
	c.ExpiresIn = time.Now().Add(time.Second * time.Duration(tokens.AuthenticationResult.ExpiresIn))

	return nil
}

func (c *CognitoToken) AutoLogin(clientId string) error {
	if time.Now().Before(c.ExpiresIn) {
		log.Println("Access is still valid")
		return nil
	}

	if err := c.RefreshingSessionToken(clientId); err != nil {
		return err
	}

	return nil
}