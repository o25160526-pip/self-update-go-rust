// manifest.rs — parse manifest-go.json (format §7)
// Testable trên Linux (không cần Tauri)

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Manifest {
    pub version: String,
    pub channel: String,
    pub published_at: String,
    pub release_notes: String,
    pub min_supported_version: String,
    pub platform: String,
    pub url: String,
    pub sha256: String,
    pub signature: String,
    pub size: u64,
    pub mandatory: bool,
}

impl Manifest {
    pub fn parse(json: &str) -> Result<Self, serde_json::Error> {
        serde_json::from_str(json)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_valid_manifest() {
        let json = r#"{
            "version": "1.0.1",
            "channel": "stable",
            "publishedAt": "2026-07-29T00:00:00Z",
            "releaseNotes": "Fix update flow",
            "minSupportedVersion": "1.0.0",
            "platform": "windows-x86_64",
            "url": "https://example.com/app.exe",
            "sha256": "abc123",
            "signature": "sig",
            "size": 12345678,
            "mandatory": false
        }"#;
        let manifest = Manifest::parse(json).unwrap();
        assert_eq!(manifest.version, "1.0.1");
        assert_eq!(manifest.channel, "stable");
        assert_eq!(manifest.platform, "windows-x86_64");
        assert_eq!(manifest.size, 12345678);
    }

    #[test]
    fn test_parse_invalid_json() {
        assert!(Manifest::parse("not json").is_err());
    }
}
