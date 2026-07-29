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
    pub fn validate(&self) -> Result<(), String> {
        semver::Version::parse(self.version.trim_start_matches('v')).map_err(|e| e.to_string())?;
        if self.channel.is_empty() || self.platform.is_empty() || self.size == 0 {
            return Err("missing required manifest field".into());
        }
        if !self.url.starts_with("https://") {
            return Err("artifact URL must use HTTPS".into());
        }
        if self.sha256.len() != 64 || !self.sha256.bytes().all(|b| b.is_ascii_hexdigit()) {
            return Err("sha256 must be 64 hex characters".into());
        }
        if self.signature.is_empty() {
            return Err("signature is required".into());
        };
        Ok(())
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    fn valid() -> Manifest {
        Manifest {
            version: "1.0.1".into(),
            channel: "stable".into(),
            published_at: "2026-07-29T00:00:00Z".into(),
            release_notes: "Fix".into(),
            min_supported_version: "1.0.0".into(),
            platform: "windows-x86_64".into(),
            url: "https://example.com/a.exe".into(),
            sha256: "a".repeat(64),
            signature: "sig".into(),
            size: 1,
            mandatory: false,
        }
    }
    #[test]
    fn validates() {
        valid().validate().unwrap()
    }
    #[test]
    fn rejects_http() {
        let mut m = valid();
        m.url = "http://example.com".into();
        assert!(m.validate().is_err())
    }
    #[test]
    fn invalid_json() {
        assert!(Manifest::parse("no").is_err())
    }
}
