use sha2::{Digest, Sha256};

/// Tính SHA-256 của data
pub fn sha256(data: &[u8]) -> String {
    let mut hasher = Sha256::new();
    hasher.update(data);
    hex::encode(hasher.finalize())
}

/// Verify SHA-256: so sánh hash tính được với expected
pub fn verify_sha256(data: &[u8], expected: &str) -> bool {
    let actual = sha256(data);
    // Constant-time comparison
    if actual.len() != expected.len() {
        return false;
    }
    let mut diff = 0u8;
    for (a, b) in actual.bytes().zip(expected.bytes()) {
        diff |= a ^ b;
    }
    diff == 0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sha256_empty() {
        // SHA-256 of empty string
        let expected = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";
        assert_eq!(sha256(b""), expected);
    }

    #[test]
    fn test_verify_sha256_match() {
        let data = b"hello world";
        let hash = sha256(data);
        assert!(verify_sha256(data, &hash));
    }

    #[test]
    fn test_verify_sha256_mismatch() {
        let data = b"hello world";
        assert!(!verify_sha256(data, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"));
    }

    #[test]
    fn test_verify_sha256_tampered_data() {
        let original = b"hello world";
        let tampered = b"hello World"; // uppercase W
        let hash = sha256(original);
        assert!(!verify_sha256(tampered, &hash));
    }
}
