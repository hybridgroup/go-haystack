package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"text/template"
)

// saveKeys writes one block of three lines for each key, in the order that the
// beacon uses them.
func saveKeys(name string, privs []string, pubs []string, hashes []string) error {
	f, err := os.Create(name + ".keys")
	if err != nil {
		return err
	}

	defer f.Close()

	for i := range privs {
		if i > 0 {
			if _, err := f.Write([]byte("\n")); err != nil {
				return err
			}
		}
		block := fmt.Sprintf("Private key: %s\nAdvertisement key: %s\nHashed adv key: %s\n",
			privs[i], pubs[i], hashes[i])
		if _, err := f.Write([]byte(block)); err != nil {
			return err
		}
	}

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
        "additionalKeys": [{{range $i, $k := .AdditionalKeys}}{{if $i}},{{end}}
            "{{$k}}"{{end}}{{if .AdditionalKeys}}
        {{end}}]
    }
]
`

// deviceData holds the values that deviceTemplate needs.
type deviceData struct {
	ID             string
	Name           string
	PrivateKey     string
	AdditionalKeys []string
}

// saveDevice writes the JSON file for macless-haystack. The first private key
// is the main key and the others go into additionalKeys, which macless-haystack
// also fetches reports for.
func saveDevice(name string, privs []string) error {
	t, err := template.New("device").Parse(deviceTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(name + ".json")
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, deviceData{
		ID:             randomInt(1000, 999999),
		Name:           name,
		PrivateKey:     privs[0],
		AdditionalKeys: privs[1:],
	})
}

// Returns an int >= min, < max
func randomInt(min, max int) string {
	return strconv.Itoa(min + rand.Intn(max-min))
}
