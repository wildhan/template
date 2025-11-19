package authorization

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	secretKey paseto.V4AsymmetricSecretKey
	publicKey paseto.V4AsymmetricPublicKey
}

func NewPasetoMaker() (*PasetoMaker, error) {

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

	return &PasetoMaker{
		secretKey: sk,
		publicKey: pk,
	}, nil
}

func (m *PasetoMaker) GenerateToken(userId string) (string, error) {
	durationMinutes := os.Getenv("TOKEN_EXPIRED_MINUTES")
	minutes, err := time.ParseDuration(durationMinutes + "m")
	if err != nil {
		return "", fmt.Errorf("invalid token duration: %w", err)
	}
	token := paseto.NewToken()
	token.SetSubject(userId)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(minutes))

	signedToken := token.V4Sign(m.secretKey, nil)

	sendingToken := strings.Split(signedToken, ".")

	return sendingToken[2], nil
}

func (m *PasetoMaker) DecryptToken(tokenString string) (*paseto.Token, error) {
	parser := paseto.NewParser()

	token := fmt.Sprintf("v4.public.%s", tokenString)
	parsed, err := parser.ParseV4Public(m.publicKey, token, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	return parsed, nil
}
