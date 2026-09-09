package haystack

import "runtime/debug"

// Version is the release version of go-haystack.
const Version = "0.1.0"

// VersionString returns the version with the VCS revision of the build, if the
// binary has one.
func VersionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}

	var rev string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}

	if rev == "" {
		return Version
	}
	if len(rev) > 7 {
		rev = rev[:7]
	}
	if dirty {
		rev += "-dirty"
	}

	return Version + " (" + rev + ")"
}
