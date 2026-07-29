use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "camelCase")]
pub struct PersistentState {
    pub current: String,
    pub pending: Option<String>,
    pub last_known_good: String,
    pub updated_from: Option<String>,
    pub rollback_attempts: u32,
    pub status: String,
}
impl PersistentState {
    pub fn new(version: &str) -> Self {
        Self {
            current: version.into(),
            pending: None,
            last_known_good: version.into(),
            updated_from: None,
            rollback_attempts: 0,
            status: "up-to-date".into(),
        }
    }
    pub fn stage(&mut self, version: &str) -> Result<(), &'static str> {
        if version.is_empty() {
            return Err("pending version is empty");
        };
        self.updated_from = Some(self.current.clone());
        self.pending = Some(version.into());
        self.status = "installing".into();
        Ok(())
    }
    pub fn mark_healthy(&mut self) -> Result<(), &'static str> {
        let v = self.pending.take().ok_or("no pending version")?;
        self.current = v.clone();
        self.last_known_good = v;
        self.rollback_attempts = 0;
        self.status = "up-to-date".into();
        Ok(())
    }
    pub fn rollback(&mut self, max: u32) -> Result<(), &'static str> {
        if self.rollback_attempts >= max {
            return Err("maximum rollback attempts reached");
        };
        if self.last_known_good.is_empty() {
            return Err("last-known-good is empty");
        };
        self.current = self.last_known_good.clone();
        self.pending = None;
        self.rollback_attempts += 1;
        self.status = "rolled-back".into();
        Ok(())
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn lifecycle() {
        let mut s = PersistentState::new("1.0.0");
        s.stage("1.0.1").unwrap();
        s.mark_healthy().unwrap();
        assert_eq!(s.current, "1.0.1");
        assert_eq!(s.last_known_good, "1.0.1");
    }
    #[test]
    fn rollback_is_limited() {
        let mut s = PersistentState::new("1.0.0");
        s.current = "1.0.1".into();
        s.rollback(1).unwrap();
        assert!(s.rollback(1).is_err());
    }
}
