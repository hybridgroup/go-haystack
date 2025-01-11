package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"text/template"
)

func saveKeys(name string, priv string, pub string, hash string) error {
	f, err := os.Create(name + ".keys")
	if err != nil {
		return err
	}

	defer f.Close()

	f.Write([]byte(fmt.Sprintf("Private key: %s\n", priv)))
	f.Write([]byte(fmt.Sprintf("Advertisement key: %s\n", pub)))
	f.Write([]byte(fmt.Sprintf("Hashed adv key: %s\n", hash)))

	return nil
}

const deviceTemplate = `[
    {
        "id": {{.ID}},
        "colorComponents": [
            0,
            1,
            0,
            1
        ],
        "name": "{{.Name}}",
        "privateKey": "{{.PrivateKey}}",
        "icon": "",
        "isDeployed": true,
        "colorSpaceName": "kCGColorSpaceExtendedSRGB",
        "usesDerivation": false,
        "isActive": false,
        "additionalKeys": []
    }
]
`

func saveDevice(name string, priv string) error {
	t, err := template.New("device").Parse(deviceTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(name + ".json")
	if err != nil {
		return err
	}

	defer f.Close()

	err = t.Execute(f, map[string]string{
		"ID":         randomInt(1000, 999999),
		"Name":       name,
		"PrivateKey": priv,
	})
	if err != nil {
		return err
	}

	return nil
}

// Returns an int >= min, < max
func randomInt(min, max int) string {
	return strconv.Itoa(min + rand.Intn(max-min))
}
