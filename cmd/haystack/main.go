package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	verboseFlag := flag.Bool("v", false, "enable verbose mode")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("either 'keys' 'flash' or 'scan' subcommand is required.", len(args))
		return
	}

	switch args[0] {
	case "keys":
		if len(args) < 2 {
			fmt.Println("Please provide a device name")
			return
		}
		if err := generateKeys(args[1], verboseFlag); err != nil {
			fmt.Println("failed to generate keys:", err)
		}
	case "flash":
		if len(args) < 3 {
			fmt.Println("Please provide a device target and name")
			return
		}
		if err := flashDevice(args[1], args[2], verboseFlag); err != nil {
			fmt.Println("failed to flash device:", err)
		}
	case "scan":
		if err := scanDevices(verboseFlag); err != nil {
			fmt.Println("failed to scan devices:", err)
		}
	default:
		fmt.Println("either 'keys' or 'flash' subcommand is required.")
		return
	}
}

func generateKeys(name string, verboseFlag *bool) error {
	// TODO: check if overwriting keys

	priv, pub, hash, err := generateKey()
	if err != nil {
		return err
	}

	// Print the keys and hash
	if *verboseFlag {
		fmt.Printf("Private key: %s\n", priv)
		fmt.Printf("Advertisement key: %s\n", pub)
		fmt.Printf("Hashed adv key: %s\n", hash)
	}

	// save keys file
	if err := saveKeys(name, priv, pub, hash); err != nil {
		return err
	}

	// save device file
	return saveDevice(name, priv)
}

func flashDevice(name string, target string, verboseFlag *bool) error {
	key, err := readKey(name)
	if err != nil {
		return err
	}

	pwd := os.Getenv("PWD")
	pth := filepath.Join(pwd, "firmware")
	if err := os.Chdir(pth); err != nil {
		panic(err)
	}
	defer os.Chdir(pwd)

	keyVal := fmt.Sprintf("-X main.AdvertisingKey='%s'", key)
	if *verboseFlag {
		fmt.Println("tinygo", "flash", "-target", target, "-ldflags", keyVal, ".")
	}

	cmd := exec.Command("tinygo", "flash", "-target", target, "-ldflags", keyVal, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readKey(name string) (string, error) {
	f, err := os.Open(name + ".keys")
	if err != nil {
		return "", err
	}
	defer f.Close()

	b := make([]byte, 1024)
	n, err := f.Read(b)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(b[:n]), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Advertisement key") {
			s := strings.Split(line, ":")
			return strings.TrimLeft(s[1], " "), nil
		}
	}

	return "", errors.New("key not found")
}
