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

// keySetTries is the number of times generateKeySet asks for a key that has a
// usable hash. About half of all keys hold a '/' in the hash.
const keySetTries = 100

// generateKeySet returns count keys as three lists, the private keys, the
// advertisement keys and the hashed advertisement keys.
func generateKeySet(count int) ([]string, []string, []string, error) {
	if count < 1 {
		return nil, nil, nil, errors.New("Key count must be 1 or more")
	}

	privs := make([]string, 0, count)
	pubs := make([]string, 0, count)
	hashes := make([]string, 0, count)

	for i := 0; i < count; i++ {
		priv, pub, hash, err := generateKeyWithRetry()
		if err != nil {
			return nil, nil, nil, err
		}
		privs = append(privs, priv)
		pubs = append(pubs, pub)
		hashes = append(hashes, hash)
	}

	return privs, pubs, hashes, nil
}

// generateKeyWithRetry returns a key that has a hash without a '/'.
func generateKeyWithRetry() (string, string, string, error) {
	for i := 0; i < keySetTries; i++ {
		priv, pub, hash, err := generateKey()
		if err == errInvalidHash {
			continue
		}
		if err != nil {
			return "", "", "", err
		}
		return priv, pub, hash, nil
	}

	return "", "", "", errInvalidHash
}
