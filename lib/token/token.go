package token

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	secretKey paseto.V4AsymmetricSecretKey
	publicKey paseto.V4AsymmetricPublicKey
}

func NewPasetoMaker(secretKey string) (*PasetoMaker, error) {

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

func (m *PasetoMaker) GenerateToken(userId string, duration time.Duration) (string, error) {
	token := paseto.NewToken()
	token.SetSubject(userId)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	signedToken := token.V4Sign(m.secretKey, nil)

	sendingToken := strings.Split(signedToken, ".")

	return sendingToken[2], nil
}
