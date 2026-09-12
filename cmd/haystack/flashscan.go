package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// scanOptions holds the build settings for the scanner.
type scanOptions struct {
	verbose  bool
	onlyMine bool
}

// scannerStackSize is the stack that the TinyScan firmware needs.
const scannerStackSize = "8kb"

// flashScanner builds the TinyScan firmware with the keys of the given devices
// and flashes it to the target.
func flashScanner(target string, names []string, opts scanOptions) error {
	keys := make([][]string, 0, len(names))
	for _, name := range names {
		if strings.ContainsAny(name, "=,;") {
			return fmt.Errorf("bad device name %q, it must have no '=' ',' or ';'", name)
		}

		k, err := readKeys(name)
		if err != nil {
			return err
		}
		keys = append(keys, k)
	}

	pwd := os.Getenv("PWD")
	pth := filepath.Join(pwd, "tinyscan")
	if err := os.Chdir(pth); err != nil {
		return err
	}
	defer os.Chdir(pwd)

	args := []string{"flash", "-target", target, "-stack-size", scannerStackSize}
	if flags := scannerLDFlags(names, keys, opts.onlyMine); flags != "" {
		args = append(args, "-ldflags", flags)
	}
	args = append(args, ".")

	if opts.verbose {
		fmt.Println("tinygo", strings.Join(args, " "))
	}

	cmd := exec.Command("tinygo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// scannerLDFlags returns the linker flags that give the scanner the devices of
// the owner. It returns an empty string if there is nothing to set.
func scannerLDFlags(names []string, keys [][]string, onlyMine bool) string {
	var flags []string

	if len(names) > 0 {
		devices := make([]string, 0, len(names))
		for i, name := range names {
			devices = append(devices, name+"="+strings.Join(keys[i], ","))
		}
		flags = append(flags, fmt.Sprintf("-X main.MyDevices='%s'", strings.Join(devices, ";")))
	}

	if onlyMine {
		flags = append(flags, "-X main.OnlyMine=true")
	}

	return strings.Join(flags, " ")
}
