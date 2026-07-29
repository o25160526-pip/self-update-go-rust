package updater

import "github.com/Masterminds/semver/v3"

// ParseVersion parse version string (hỗ trợ 'v' prefix).
func ParseVersion(v string) (*semver.Version, error) {
	v = stripVPrefix(v)
	return semver.NewVersion(v)
}

// IsNewer trả về true nếu latest > current.
func IsNewer(current, latest string) (bool, error) {
	cur, err := ParseVersion(current)
	if err != nil {
		return false, err
	}
	lat, err := ParseVersion(latest)
	if err != nil {
		return false, err
	}
	return lat.GreaterThan(cur), nil
}

// MeetsMinimum trả về true nếu latest >= minSupported.
func MeetsMinimum(latest, minSupported string) (bool, error) {
	lat, err := ParseVersion(latest)
	if err != nil {
		return false, err
	}
	min, err := ParseVersion(minSupported)
	if err != nil {
		return false, err
	}
	return !lat.LessThan(min), nil
}

func stripVPrefix(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}
	return v
}
