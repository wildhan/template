package paseto

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"template/config/authorization"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoImpl struct {
	secretKey    *paseto.V4AsymmetricSecretKey
	publicKey    *paseto.V4AsymmetricPublicKey
	symmetricKey *paseto.V4SymmetricKey
}

func NewPasetoImpl() (authorization.AuthToken, error) {
	secretKey := os.Getenv("SECRET_KEY_64_BYTES")
	raw, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key encoding: %w", err)
	}

	sk, err := paseto.NewV4AsymmetricSecretKeyFromBytes(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key size: %w", err)
	}
	pk := sk.Public()
	fmt.Printf("Public Key: %s", base64.StdEncoding.EncodeToString(pk.ExportBytes()))

	return &PasetoImpl{
		secretKey:    &sk,
		publicKey:    &pk,
		symmetricKey: nil,
	}, nil
}

func (p *PasetoImpl) GenerateToken(tc authorization.TokenContainer) (string, error) {
	durationMinutes := os.Getenv("TOKEN_EXPIRED_MINUTES")
	minutes, err := time.ParseDuration(durationMinutes + "m")
	if err != nil {
		return "", fmt.Errorf("invalid token duration: %w", err)
	}
	token := paseto.NewToken()
	token.SetSubject(tc.Sender)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(minutes))

	signedToken := token.V4Sign(*p.secretKey, nil)

	sendingToken := strings.Split(signedToken, ".")

	return sendingToken[2], nil
}

func (p *PasetoImpl) ValidateToken(tokenString string) (*authorization.TokenContainer, error) {
	parser := paseto.NewParser()

	token := fmt.Sprintf("v4.public.%s", tokenString)
	parsed, err := parser.ParseV4Public(*p.publicKey, token, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	// Check Expiration Time
	expiredAt, err := parsed.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("invalid token expiration: %w", err)
	}

	if expiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	sub, err := parsed.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("invalid token subject: %w", err)
	}

	tc := &authorization.TokenContainer{
		Sender: sub,
	}
	return tc, nil
}

func (p *PasetoImpl) GenerateRefreshToken(rtc authorization.RefreshTokenContainer) (string, error) {
	return "", nil
}

func (p *PasetoImpl) ValidateRefreshToken(rtString string) (*authorization.RefreshTokenContainer, error) {
	return nil, nil
}
