package updater

import "encoding/json"

// Manifest đại diện cho manifest format theo §7 của prompt.
type Manifest struct {
	Version             string `json:"version"`
	Channel             string `json:"channel"`
	PublishedAt         string `json:"publishedAt"`
	ReleaseNotes        string `json:"releaseNotes"`
	MinSupportedVersion string `json:"minSupportedVersion"`
	Platform            string `json:"platform"`
	URL                 string `json:"url"`
	SHA256              string `json:"sha256"`
	Signature           string `json:"signature"`
	Size                int64  `json:"size"`
	Mandatory           bool   `json:"mandatory"`
}

// ParseManifest parse manifest JSON.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
