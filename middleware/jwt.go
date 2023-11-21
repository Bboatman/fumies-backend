package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Nickname    *string `json:"nickname"`
	Name        *string `json:"name"`
	PictureUrl  *string `json:"picture"`
	Scope       *string `json:"scope"`
	UpdatedAt   *string `json:"updated_at"`
	UserSubject *string `json:"sub"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}
func JWTValidator() gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("authorization")
		if authHeader == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Failed to validate JWT."))
			return
		}

		authTokenSlice := strings.Split(authHeader, " ")
		if len(authTokenSlice) < 2 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Failed to validate JWT."))
			return
		}

		validationResult, err := jwtValidator.ValidateToken(c, authTokenSlice[1])
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Failed to validate JWT."))
			return
		}
		claims := validationResult.(*validator.ValidatedClaims)

		customClaims := claims.CustomClaims.(*CustomClaims)
		if customClaims != nil && customClaims.UserSubject != nil {
			c.Set("user_id", *customClaims.UserSubject)
		} else {
			c.Set("user_id", claims.RegisteredClaims.Subject)
		}

		// before request
		c.Next()
		// after request
	}
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	if c.Scope == nil {
		return false
	}

	result := strings.Split(*c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
