package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"book-store/internal/config"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type Auth struct {
	jwks     keyfunc.Keyfunc
	issuer   string
	audience string
}

func NewAuth(lc fx.Lifecycle, cfg *config.Config) (*Auth, error) {
	if cfg.AzureTenantID == "" {
		return nil, errors.New("AZURE_TENANT_ID is required")
	}
	if cfg.AzureClientID == "" {
		return nil, errors.New("AZURE_CLIENT_ID is required")
	}

	jwksURL := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/discovery/v2.0/keys",
		cfg.AzureTenantID,
	)

	ctx, cancel := context.WithCancel(context.Background())

	k, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})

	return &Auth{
		jwks:     k,
		issuer:   fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", cfg.AzureTenantID),
		audience: cfg.AzureClientID,
	}, nil
}

func (a *Auth) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, a.jwks.Keyfunc,
			jwt.WithIssuer(a.issuer),
			jwt.WithAudience(a.audience),
			jwt.WithValidMethods([]string{"RS256"}),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if oid, ok := claims["oid"].(string); ok {
				c.Set("user_id", oid)
			}
		}

		c.Next()
	}
}
