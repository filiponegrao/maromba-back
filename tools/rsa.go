package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	rsaKeySize = 2048
)

func hash(data []byte) []byte {
	s := sha256.Sum256(data)
	return s[:]
}

// MARK: Public and Private Keys

func generateKeyBytes() (privateBytes, publicBytes []byte, err error) {
	pri, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, nil, err
	}
	pub := &pri.PublicKey
	priBytes := x509.MarshalPKCS1PrivateKey(pri)
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	return priBytes, pubBytes, nil
}

// GenerateKeys : GenerateKeys
func GenerateRSAKeys() (pri64, pub64 string, err error) {
	pri, pub, err := generateKeyBytes()
	if err != nil {
		return "", "", nil
	}
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pri,
		},
	)
	privString := string(privkeyPem)
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pub,
		},
	)
	pubString := string(pubkeyPem)
	return privString, pubString, nil
}

func PrivateRSAKeyFromString(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func PublicRSAKeyFromString(key string, pkcs1 bool) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	if pkcs1 {
		// PKCS1
		pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pub := &pri.PublicKey
		return pub, nil

	} else {
		// PKCS8
		pri, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		p, ok := pri.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("invalid public key")
		}
		return p, nil
	}
}

// MARK: Encrypt & Decrypt

// PrivateDecrypt : PrivateDecrypt
func DecryptRSAWithPrivateKey(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, data)
}

func PublicVerify(key *rsa.PublicKey, data []byte, signature []byte) error {
	hashed := hash(data)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed, signature)
}
