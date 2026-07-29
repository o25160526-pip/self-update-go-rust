// policy.rs — parse update-policy.json
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct UpdatePolicy {
    pub channel: String,
    pub auto_check_on_startup: bool,
    pub auto_check_interval_minutes: u32,
    pub auto_download: bool,
    pub auto_install: bool,
    pub require_user_confirmation: bool,
    pub restart_automatically: bool,
    pub allow_downgrade: bool,
    pub rollback_on_startup_failure: bool,
    pub health_check_timeout_seconds: u32,
    pub max_rollback_attempts: u32,
    pub min_supported_version: String,
}

impl Default for UpdatePolicy {
    fn default() -> Self {
        Self {
            channel: "stable".to_string(),
            auto_check_on_startup: true,
            auto_check_interval_minutes: 60,
            auto_download: true,
            auto_install: true,
            require_user_confirmation: false,
            restart_automatically: true,
            allow_downgrade: false,
            rollback_on_startup_failure: true,
            health_check_timeout_seconds: 30,
            max_rollback_attempts: 1,
            min_supported_version: "1.0.0".to_string(),
        }
    }
}

impl UpdatePolicy {
    pub fn parse(json: &str) -> Result<Self, serde_json::Error> {
        serde_json::from_str(json)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_policy() {
        let p = UpdatePolicy::default();
        assert!(p.auto_check_on_startup);
        assert!(!p.require_user_confirmation);
        assert!(p.rollback_on_startup_failure);
        assert_eq!(p.max_rollback_attempts, 1);
    }

    #[test]
    fn test_parse_policy() {
        let json = r#"{
            "channel": "stable",
            "autoCheckOnStartup": true,
            "autoCheckIntervalMinutes": 30,
            "autoDownload": true,
            "autoInstall": true,
            "requireUserConfirmation": false,
            "restartAutomatically": true,
            "allowDowngrade": false,
            "rollbackOnStartupFailure": true,
            "healthCheckTimeoutSeconds": 30,
            "maxRollbackAttempts": 1,
            "minSupportedVersion": "1.0.0"
        }"#;
        let p = UpdatePolicy::parse(json).unwrap();
        assert_eq!(p.auto_check_interval_minutes, 30);
    }
}
