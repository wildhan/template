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

	symmByte := os.Getenv("SYMMETRIC_KEY_32_BYTES")
	skRaw, err := base64.StdEncoding.DecodeString(symmByte)
	if err != nil {
		return nil, fmt.Errorf("invalid symmetric key encoding: %w", err)
	}

	syKey, err := paseto.V4SymmetricKeyFromBytes(skRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid symmetric key size: %w", err)
	}

	return &PasetoImpl{
		secretKey:    &sk,
		publicKey:    &pk,
		symmetricKey: &syKey,
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
	durationHours := os.Getenv("REFRESH_TOKEN_EXPIRED_HOURS")
	hours, err := time.ParseDuration(durationHours + "h")
	if err != nil {
		return "", fmt.Errorf("invalid token duration: %w", err)
	}
	token := paseto.NewToken()
	token.SetSubject(rtc.TC.Sender)
	token.SetString("email", rtc.Email)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(hours))

	signedToken := token.V4Encrypt(*p.symmetricKey, nil)

	sendingToken := strings.Split(signedToken, ".")

	return sendingToken[2], nil
}

func (p *PasetoImpl) ValidateRefreshToken(rtString string) (*authorization.RefreshTokenContainer, error) {
	token := fmt.Sprintf("v4.local.%s", rtString)
	parser := paseto.NewParser()
	parsed, err := parser.ParseV4Local(*p.symmetricKey, token, nil)
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

	email, err := parsed.GetString("email")
	if err != nil {
		return nil, fmt.Errorf("invalid token email: %w", err)
	}
	rtc := &authorization.RefreshTokenContainer{
		TC: authorization.TokenContainer{
			Sender: sub,
		},
		Email: email,
	}
	return rtc, nil
}
