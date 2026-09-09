package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestGenerateKeySet checks the count, the content and the hashes of a set.
func TestGenerateKeySet(t *testing.T) {
	for _, count := range []int{1, 3, 12} {
		privs, pubs, hashes, err := generateKeySet(count)
		if err != nil {
			t.Errorf("count %d: %v", count, err)
			continue
		}
		if len(privs) != count || len(pubs) != count || len(hashes) != count {
			t.Errorf("count %d: got %d %d %d keys", count, len(privs), len(pubs), len(hashes))
			continue
		}

		seen := make(map[string]bool, count)
		for i := range privs {
			if seen[privs[i]] {
				t.Errorf("count %d: key %d is a repeat", count, i)
			}
			seen[privs[i]] = true

			// macless-haystack puts the hash in a URL path, so a '/' in it
			// stops the reports.
			if strings.Contains(hashes[i], "/") {
				t.Errorf("count %d: hash %d holds a '/'", count, i)
			}

			val, err := base64.StdEncoding.DecodeString(pubs[i])
			if err != nil || len(val) != 28 {
				t.Errorf("count %d: adv key %d is %d bytes, %v", count, i, len(val), err)
			}

			// A key with a leading zero must keep the full size.
			val, err = base64.StdEncoding.DecodeString(privs[i])
			if err != nil || len(val) != 28 {
				t.Errorf("count %d: private key %d is %d bytes, %v", count, i, len(val), err)
			}
		}
	}

	if _, _, _, err := generateKeySet(0); err == nil {
		t.Error("count 0: got no error, want one")
	}
}

// TestSaveAndReadKeys checks that the file keeps the keys in order.
func TestSaveAndReadKeys(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"one", 1},
		{"twelve", 12},
	}

	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	for _, tc := range tests {
		privs, pubs, hashes, err := generateKeySet(tc.count)
		if err != nil {
			t.Fatal(err)
		}
		if err := saveKeys(tc.name, privs, pubs, hashes); err != nil {
			t.Fatal(err)
		}

		got, err := readKeys(tc.name)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		if len(got) != tc.count {
			t.Errorf("%s: got %d keys, want %d", tc.name, len(got), tc.count)
			continue
		}
		for i := range got {
			if got[i] != pubs[i] {
				t.Errorf("%s: key %d is %q, want %q", tc.name, i, got[i], pubs[i])
			}
		}
	}

	if _, err := readKeys("missing"); err == nil {
		t.Error("missing file: got no error, want one")
	}

	if err := os.WriteFile("empty.keys", []byte("nothing here\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readKeys("empty"); err == nil {
		t.Error("file without a key: got no error, want one")
	}
}

// TestSaveDevice checks the JSON file that macless-haystack reads.
func TestSaveDevice(t *testing.T) {
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	tests := []struct {
		name  string
		privs []string
	}{
		{"single", []string{"AAA"}},
		{"set", []string{"AAA", "BBB", "CCC"}},
	}

	for _, tc := range tests {
		if err := saveDevice(tc.name, tc.privs); err != nil {
			t.Fatal(err)
		}

		b, err := os.ReadFile(tc.name + ".json")
		if err != nil {
			t.Fatal(err)
		}

		var devices []struct {
			Name           string   `json:"name"`
			PrivateKey     string   `json:"privateKey"`
			AdditionalKeys []string `json:"additionalKeys"`
		}
		if err := json.Unmarshal(b, &devices); err != nil {
			t.Errorf("%s: %v\n%s", tc.name, err, b)
			continue
		}
		if len(devices) != 1 {
			t.Errorf("%s: got %d devices, want 1", tc.name, len(devices))
			continue
		}

		dev := devices[0]
		if dev.Name != tc.name {
			t.Errorf("%s: name is %q", tc.name, dev.Name)
		}
		if dev.PrivateKey != tc.privs[0] {
			t.Errorf("%s: privateKey is %q, want %q", tc.name, dev.PrivateKey, tc.privs[0])
		}
		if len(dev.AdditionalKeys) != len(tc.privs)-1 {
			t.Errorf("%s: got %d additional keys, want %d", tc.name, len(dev.AdditionalKeys), len(tc.privs)-1)
			continue
		}
		for i, key := range dev.AdditionalKeys {
			if key != tc.privs[i+1] {
				t.Errorf("%s: additional key %d is %q, want %q", tc.name, i, key, tc.privs[i+1])
			}
		}
	}
}
