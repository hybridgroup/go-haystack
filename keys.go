package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

var errInvalidHash = errors.New("Hash contains '/'")

func generateKey() (string, string, string, error) {
	// Generate ECDSA private key using P-224 curve
	pk, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return "", "", "", err
	}

	// Extract the raw private key bytes
	privateKeyBytes := pk.D.Bytes()

	// Ensure the private key is 28 bytes long (P-224 curve)
	if len(privateKeyBytes) != 28 {
		return "", "", "", errors.New("Private key is not 28 bytes long")
	}

	// Encode the raw private key to Base64
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyBytes)

	// extract raw public key bytes
	publicKeyBytes := pk.PublicKey.X.Bytes()

	// Encode the public key to Base64
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)

	// Hash the public key using SHA-256
	hash := sha256.Sum256(publicKeyBytes)
	hashBase64 := base64.StdEncoding.EncodeToString(hash[:])

	// make sure not '/' in the base64 string
	if strings.Contains(hashBase64, "/") {
		return "", "", "", errInvalidHash
	}

	return privateKeyBase64, publicKeyBase64, hashBase64, nil
}
