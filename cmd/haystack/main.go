package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	verboseFlag := flag.Bool("v", false, "enable verbose mode")
	batteryFlag := flag.Bool("battery", false, "build for a device on a battery, which uses as little current as possible")
	txPowerFlag := flag.String("txpower", "", "radio transmit power in dBm, for example -8. Empty keeps the default power")
	dcdc0Flag := flag.Bool("dcdc0", false, "turn the DC/DC converter of the VDDH stage on. Needs a board powered through VDDH, and often gains nothing from a battery")
	batteryPinFlag := flag.String("batterypin", "", "GPIO number that reads a battery divider, for example 2. Only an ESP32-C3 or ESP32-S3 board needs it")
	batteryDividerFlag := flag.String("batterydivider", "", "ratio of the battery divider, for example 2/1. A single number is a ratio to 1")
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
		opts := flashOptions{
			verbose:        *verboseFlag,
			battery:        *batteryFlag,
			dcdc0:          *dcdc0Flag,
			txPower:        *txPowerFlag,
			batteryPin:     *batteryPinFlag,
			batteryDivider: *batteryDividerFlag,
		}
		if err := flashDevice(args[1], args[2], opts); err != nil {
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

// flashOptions holds the build settings that the flags give.
type flashOptions struct {
	verbose        bool
	battery        bool
	dcdc0          bool
	txPower        string
	batteryPin     string
	batteryDivider string
}

// espTargets are the TinyGo targets that use the radio in an ESP32-C3 or
// ESP32-S3 chip.
var espTargets = []string{
	"xiao-esp32c3",
	"xiao-esp32s3",
	"esp32c3-supermini",
	"esp32s3-supermini",
	"esp32c3-generic",
	"esp32s3-generic",
	"qtpy-esp32c3",
	"m5stamp-c3",
}

// isESPTarget reports if the target uses the radio in an ESP32 chip.
func isESPTarget(target string) bool {
	return slices.Contains(espTargets, target)
}

func flashDevice(name string, target string, opts flashOptions) error {
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

	esp := isESPTarget(target)

	keyVal := fmt.Sprintf("-X main.AdvertisingKey='%s'", key)
	if opts.txPower != "" {
		keyVal += fmt.Sprintf(" -X main.TxPower=%s", opts.txPower)
		if esp {
			fmt.Println("note: this target does not support the transmit power setting, so the radio keeps its default power")
		}
	}
	if opts.dcdc0 {
		keyVal += " -X main.DCDC0=on"
	}
	if opts.batteryPin != "" {
		keyVal += fmt.Sprintf(" -X main.BatteryPin=%s", opts.batteryPin)
	}
	if opts.batteryDivider != "" {
		keyVal += fmt.Sprintf(" -X main.BatteryDivider=%s", opts.batteryDivider)
	}

	args := []string{"flash", "-target", target}
	if opts.battery {
		if esp {
			// The console of an ESP32-C3 or ESP32-S3 is part of the USB block,
			// which stays on, so this only removes the messages.
			fmt.Println("note: this target keeps its serial port, because turning it off saves almost no current and removes all messages")
		} else {
			// The USB peripheral uses current for no purpose on a battery.
			args = append(args, "-serial=none")
		}
	}
	args = append(args, "-ldflags", keyVal, ".")

	if opts.verbose {
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
