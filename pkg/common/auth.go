package common

import (
	"fmt"
	"sachahjkl/htmx_go/pkg/model"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func UserFromJwt(db *gorm.DB, userJWT *jwt.Token) (*model.User, error) {

	if userJWT == nil {
		return nil, fmt.Errorf("user JWT is nil")
	}

	claims := userJWT.Claims.(jwt.MapClaims)

	userUnknown := claims[USER_CLAIM_KEY]

	if userUnknown == nil {
		return nil, fmt.Errorf("user id not in claims")
	}

	userId, ok := userUnknown.(float64)
	if !ok {
		return nil, fmt.Errorf("couldn't convert claim to float64")
	}

	return model.GetUser(db, uint(userId))

}
