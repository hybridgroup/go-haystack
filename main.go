package main

import (
	"flag"
	"fmt"
)

func main() {
	nameFlag := flag.String("name", "", "name of the device to be added")
	verboseFlag := flag.Bool("verbose", false, "enable verbose mode")

	flag.Parse()
	if *nameFlag == "" {
		fmt.Println("Please provide a device name using the -name flag")
		return
	}

	priv, pub, hash := generateKey()

	// Print the keys and hash
	if *verboseFlag {
		fmt.Printf("Private key: %s\n", priv)
		fmt.Printf("Advertisement key: %s\n", pub)
		fmt.Printf("Hashed adv key: %s\n", hash)
	}

	// save keys file
	saveKeys(*nameFlag, priv, pub, hash)

	// save device file
	saveDevice(*nameFlag, priv)
}
