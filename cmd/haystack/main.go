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
	batteryFlag := flag.Bool("battery", false, "build for a device on a battery, which uses as little current as possible")
	txPowerFlag := flag.String("txpower", "", "radio transmit power in dBm, for example -8. Empty keeps the default power")
	dcdc0Flag := flag.Bool("dcdc0", false, "turn the DC/DC converter of the VDDH stage on. Needs a board powered through VDDH, and often gains nothing from a battery")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("subcommand required. valid subcommands are 'keys' 'flash' 'scan'")
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
			fmt.Println("Please provide a device name and target")
			return
		}
		if err := flashDevice(args[1], args[2], verboseFlag, batteryFlag, dcdc0Flag, txPowerFlag); err != nil {
			fmt.Println("failed to flash device:", err)
		}
	case "scan":
		if err := scanDevices(verboseFlag); err != nil {
			fmt.Println("failed to scan devices:", err)
		}
	default:
		fmt.Println("subcommand required. valid subcommands are 'keys' 'flash' 'scan'")
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

func flashDevice(name string, target string, verboseFlag, batteryFlag, dcdc0Flag *bool, txPowerFlag *string) error {
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
	if *txPowerFlag != "" {
		keyVal += fmt.Sprintf(" -X main.TxPower=%s", *txPowerFlag)
	}
	if *dcdc0Flag {
		keyVal += " -X main.DCDC0=on"
	}

	args := []string{"flash", "-target", target}
	if *batteryFlag {
		// The USB peripheral uses current for no purpose on a battery.
		args = append(args, "-serial=none")
	}
	args = append(args, "-ldflags", keyVal, ".")

	if *verboseFlag {
		fmt.Println("tinygo", strings.Join(args, " "))
	}

	cmd := exec.Command("tinygo", args...)
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
