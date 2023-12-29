package secp256k1identity

import (
	"crypto/sha256"
	"errors"

	"github.com/aviate-labs/agent-go/principal"
	"github.com/btcsuite/btcd/btcec"
	"github.com/dfinity/keysmith/codec"
)

type Secp256k1Identity struct {
	privateKey          *btcec.PrivateKey
	publicKey           *btcec.PublicKey
	derEncodedPublicKey []byte
}

func NewSecp256k1Identity(privateKey *btcec.PrivateKey) (*Secp256k1Identity, error) {
	publicKey := privateKey.PubKey()
	der, err := codec.EncodeECPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &Secp256k1Identity{
		privateKey:          privateKey,
		publicKey:           publicKey,
		derEncodedPublicKey: der,
	}, nil
}

func NewSecp256k1IdentityFromPEM(data []byte) (*Secp256k1Identity, error) {
	privateKey, err := codec.PEMToECPrivKey(data)
	if err != nil {
		return nil, errors.New("invalid pem data")
	}
	return NewSecp256k1Identity(privateKey)
}

func (p *Secp256k1Identity) Sender() principal.Principal {
	return principal.NewSelfAuthenticating(p.derEncodedPublicKey)
}

func (p *Secp256k1Identity) Sign(msg []byte) []byte {
	hash := sha256.New()
	hash.Write(msg)
	hashData := hash.Sum(nil)

	sig, _ := p.privateKey.Sign(hashData)
	return codec.EncodeECSig(sig)
}

func (p *Secp256k1Identity) PublicKey() []byte {
	return p.derEncodedPublicKey
}
