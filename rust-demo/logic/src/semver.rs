// semver.rs — logic so sánh version (testable trên Linux)
// Dùng crate semver (đã có qua Tauri deps)

pub fn parse(version: &str) -> Option<semver::Version> {
    // Strip leading 'v' if present
    let v = version.strip_prefix('v').unwrap_or(version);
    semver::Version::parse(v).ok()
}

/// So sánh current vs latest. Trả về true nếu latest > current.
pub fn is_newer(current: &str, latest: &str) -> Result<bool, String> {
    let cur = parse(current).ok_or(format!("invalid current version: {}", current))?;
    let lat = parse(latest).ok_or(format!("invalid latest version: {}", latest))?;
    Ok(lat > cur)
}

/// Kiểm tra downgrade: latest phải >= min_supported
pub fn meets_minimum(latest: &str, min_supported: &str) -> Result<bool, String> {
    let lat = parse(latest).ok_or(format!("invalid latest version: {}", latest))?;
    let min = parse(min_supported).ok_or(format!("invalid min version: {}", min_supported))?;
    Ok(lat >= min)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_with_v_prefix() {
        assert_eq!(parse("v1.0.0"), parse("1.0.0"));
    }

    #[test]
    fn test_is_newer() {
        assert!(is_newer("1.0.0", "1.0.1").unwrap());
        assert!(!is_newer("1.0.1", "1.0.0").unwrap());
        assert!(!is_newer("1.0.0", "1.0.0").unwrap());
    }

    #[test]
    fn test_meets_minimum() {
        assert!(meets_minimum("1.0.1", "1.0.0").unwrap());
        assert!(meets_minimum("1.0.0", "1.0.0").unwrap());
        assert!(!meets_minimum("0.9.0", "1.0.0").unwrap());
    }

    #[test]
    fn test_invalid_version() {
        assert!(parse("not-a-version").is_none());
        assert!(is_newer("1.0.0", "invalid").is_err());
    }
}
