package main

import (
	"fmt"
)

func main() {
	priv, pub, hash := generateKey()

	// Print the keys and hash
	fmt.Printf("Private Key (Base64): %s\n", priv)
	fmt.Printf("Advertising Key (Base64): %s\n", pub)
	fmt.Printf("Advertising Key Hash (Base64): %s\n", hash)
}
