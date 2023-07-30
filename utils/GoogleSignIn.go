package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/idtoken"
)

func ValidateGoogleIdToken(idToken string) (*idtoken.Payload, error) {
	ctx := context.Background()

	clientIds := []string{
		"987503717012-ve2cvppbv230g248pgvmibep4mam7gun.apps.googleusercontent.com", // Android client ID
		"987503717012-hi1a327ahjs8nb6plb3b6s8918q0r8e5.apps.googleusercontent.com", // iOS client ID
	}

	var payload *idtoken.Payload
	var err error

	for _, clientId := range clientIds {
		// Verify the ID token with the client ID.
		payload, err = idtoken.Validate(ctx, idToken, clientId)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to validate id token: %v", err)
	}

	return payload, nil
}

func GetNameFromIdToken(idToken string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse id token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to parse token claims")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return "", errors.New("name field not found in id token")
	}

	return name, nil
}
