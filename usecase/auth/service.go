package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/org39/webapp-tutorial-backend/entity"

	"github.com/dgrijalva/jwt-go"
)

type Service struct {
	Secret               string        `inject:"usecase.auth.secret"`
	AccessTokenDuration  time.Duration `inject:"usecase.auth.access_token_duration"`
	RefreshTokenDuration time.Duration `inject:"usecase.auth.refresh_token_duration"`
}

func NewService(options ...func(*Service) error) (Usecase, error) {
	u := &Service{
		Secret: "",
		// access token duration, 10 minute default
		AccessTokenDuration: 10 * time.Minute,
		// refresh token duration, 7 days default
		RefreshTokenDuration: 7 * 24 * time.Hour,
	}

	for _, option := range options {
		if err := option(u); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func WithSecret(s string) func(*Service) error {
	return func(u *Service) error {
		u.Secret = s
		return nil
	}
}

func (u *Service) GenereateToken(ctx context.Context, id string) (*entity.AuthTokenPair, error) {
	if err := entity.NewValidator().ValidateID(id); err != nil {
		return nil, fmt.Errorf("%s: invalid token request: %w", err, ErrInvalidRequest)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	now := time.Now()

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = now.Add(u.AccessTokenDuration).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret.
	t, err := token.SignedString([]byte(u.Secret))
	if err != nil {
		return nil, fmt.Errorf("%s: generate token error: %w", err, ErrSystemError)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["exp"] = now.Add(u.RefreshTokenDuration).Unix()

	rt, err := refreshToken.SignedString([]byte(u.Secret))
	if err != nil {
		return nil, fmt.Errorf("%s: generate refresh token error: %w", err, ErrSystemError)
	}

	return entity.NewFactory().NewAuthTokenPair(t, rt), nil
}

func (u *Service) RefreshToken(ctx context.Context, refreshToken string) (*entity.AuthTokenPair, error) {
	if err := entity.NewValidator().ValidateToken(refreshToken); err != nil {
		return nil, fmt.Errorf("%s: invalid refresh request: %w", err, ErrInvalidRequest)
	}

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method(%v): %w", token.Header["alg"], ErrInvalidRequest)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("secret")
		return []byte(u.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnauthorized)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %w", ErrUnauthorized)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %w", ErrUnauthorized)
	}

	newTokenPair, err := u.GenereateToken(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrUnauthorized)
	}

	return newTokenPair, nil
}

func (u *Service) VerifyToken(ctx context.Context, accessToken string) (string, error) {
	if err := entity.NewValidator().ValidateToken(accessToken); err != nil {
		return "", fmt.Errorf("%s: invalid verify request: %w", err, ErrInvalidRequest)
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method(%v): %w", token.Header["alg"], ErrInvalidRequest)
		}

		return []byte(u.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", err, ErrUnauthorized)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token: %w", ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims: %w", ErrUnauthorized)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid claims: %w", ErrUnauthorized)
	}

	return id, nil
}
