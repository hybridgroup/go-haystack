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

	"github.com/hybridgroup/go-haystack"
)

// usage names the subcommands of the tool.
const usage = "subcommand required. valid subcommands are 'keys' 'flash' 'flashscan' 'scan' 'version'"

func main() {
	verboseFlag := flag.Bool("v", false, "enable verbose mode")
	batteryFlag := flag.Bool("battery", false, "build for a device on a battery, which uses as little current as possible")
	txPowerFlag := flag.String("txpower", "", "radio transmit power in dBm, for example -8. Empty keeps the default power")
	dcdc0Flag := flag.Bool("dcdc0", false, "turn the DC/DC converter of the VDDH stage on. Needs a board powered through VDDH, and often gains nothing from a battery")
	batteryPinFlag := flag.String("batterypin", "", "GPIO number that reads a battery divider, for example 2. Only an ESP32-C3 or ESP32-S3 board needs it")
	batteryDividerFlag := flag.String("batterydivider", "", "ratio of the battery divider, for example 2/1. A single number is a ratio to 1")
	batteryTypeFlag := flag.String("batterytype", "", "cell that the beacon uses, one of lipo, cr2032, cr1220 or aa-alkaline. Empty is lipo")
	batteryThresholdsFlag := flag.String("batterythresholds", "", "full, medium and low battery voltages in millivolts, for example 2900/2750/2600. It wins over -batterytype")
	keysFlag := flag.Int("keys", defaultKeyCount, "how many keys to generate for a device, which the beacon then uses in turn")
	onlyMineFlag := flag.Bool("onlymine", false, "the scanner shows only your own devices")
	rotateFlag := flag.String("rotate", defaultKeyRotation, "how long the beacon uses each key, for example 5m. An empty value stops the rotation")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println(usage)
		return
	}

	switch args[0] {
	case "keys":
		if len(args) < 2 {
			fmt.Println("Please provide a device name")
			return
		}
		if err := generateKeys(args[1], *keysFlag, *verboseFlag); err != nil {
			fmt.Println("failed to generate keys:", err)
		}
	case "flash":
		if len(args) < 3 {
			fmt.Println("Please provide a device name and target")
			return
		}
		opts := flashOptions{
			verbose:           *verboseFlag,
			battery:           *batteryFlag,
			dcdc0:             *dcdc0Flag,
			txPower:           *txPowerFlag,
			batteryPin:        *batteryPinFlag,
			batteryDivider:    *batteryDividerFlag,
			batteryType:       *batteryTypeFlag,
			batteryThresholds: *batteryThresholdsFlag,
			rotate:            *rotateFlag,
		}
		if err := flashDevice(args[1], args[2], opts); err != nil {
			fmt.Println("failed to flash device:", err)
		}
	case "flashscan":
		if len(args) < 2 {
			fmt.Println("Please provide a target")
			return
		}
		opts := scanOptions{
			verbose:  *verboseFlag,
			onlyMine: *onlyMineFlag,
		}
		if err := flashScanner(args[1], args[2:], opts); err != nil {
			fmt.Println("failed to flash scanner:", err)
		}
	case "scan":
		if err := scanDevices(verboseFlag); err != nil {
			fmt.Println("failed to scan devices:", err)
		}
	case "version":
		fmt.Println("haystack", haystack.VersionString())
	default:
		fmt.Println(usage)
		return
	}
}

func generateKeys(name string, count int, verbose bool) error {
	// TODO: check if overwriting keys

	privs, pubs, hashes, err := generateKeySet(count)
	if err != nil {
		return err
	}

	// Print the keys and hashes
	if verbose {
		for i := range privs {
			fmt.Printf("Private key: %s\n", privs[i])
			fmt.Printf("Advertisement key: %s\n", pubs[i])
			fmt.Printf("Hashed adv key: %s\n", hashes[i])
		}
	}

	// save keys file
	if err := saveKeys(name, privs, pubs, hashes); err != nil {
		return err
	}

	// save device file
	return saveDevice(name, privs)
}

// defaultKeyCount is how many keys a device gets. With defaultKeyRotation the
// beacon uses the whole set in one hour and then starts again.
const defaultKeyCount = 12

// defaultKeyRotation is how long the beacon uses each key.
const defaultKeyRotation = "5m"

// flashOptions holds the build settings that the flags give.
type flashOptions struct {
	verbose        bool
	battery        bool
	dcdc0          bool
	txPower        string
	batteryPin     string
	batteryDivider string
	// batteryType and batteryThresholds set the voltages that give the battery
	// status. batteryThresholds wins over batteryType.
	batteryType       string
	batteryThresholds string
	rotate            string
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
	keys, err := readKeys(name)
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

	keyVal := fmt.Sprintf("-X main.AdvertisingKey='%s'", strings.Join(keys, ","))
	if len(keys) > 1 && opts.rotate != "" {
		keyVal += fmt.Sprintf(" -X main.KeyRotation=%s", opts.rotate)
	}
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
	if opts.batteryType != "" {
		keyVal += fmt.Sprintf(" -X main.BatteryType=%s", opts.batteryType)
	}
	if opts.batteryThresholds != "" {
		keyVal += fmt.Sprintf(" -X main.BatteryThresholds=%s", opts.batteryThresholds)
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

// readKeys returns every advertisement key in the file of a device, in the
// order that the beacon uses them.
func readKeys(name string) ([]string, error) {
	b, err := os.ReadFile(name + ".keys")
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, line := range strings.Split(string(b), "\n") {
		if !strings.Contains(line, "Advertisement key") {
			continue
		}
		s := strings.SplitN(line, ":", 2)
		keys = append(keys, strings.TrimSpace(s[1]))
	}

	if len(keys) == 0 {
		return nil, errors.New("key not found")
	}

	return keys, nil
}
