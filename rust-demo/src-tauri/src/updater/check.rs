use std::sync::{Mutex, OnceLock};
use tauri::AppHandle;
use tauri_plugin_updater::UpdaterExt;

static STATE: OnceLock<Mutex<String>> = OnceLock::new();
fn state_lock() -> &'static Mutex<String> {
    STATE.get_or_init(|| Mutex::new("checking".into()))
}
pub fn get_state() -> String {
    state_lock()
        .lock()
        .unwrap_or_else(|e| e.into_inner())
        .clone()
}
pub fn set_state(value: &str) {
    *state_lock().lock().unwrap_or_else(|e| e.into_inner()) = value.into();
}

pub async fn check_on_startup(app: AppHandle, offline: bool) {
    set_state("checking");
    if offline {
        let result = crate::updater::mock::check();
        set_state(if result.has_update {
            "update-available"
        } else {
            "up-to-date"
        });
        return;
    }
    let updater = match app.updater() {
        Ok(v) => v,
        Err(e) => {
            eprintln!("updater init failed: {e}");
            set_state("failed");
            return;
        }
    };
    match updater.check().await {
        Ok(Some(update)) => {
            set_state("update-available");
            if let Err(e) = crate::updater::install::download_install_restart(app, update).await {
                eprintln!("update failed: {e}");
                set_state("failed");
            }
        }
        Ok(None) => set_state("up-to-date"),
        Err(e) => {
            eprintln!("update check failed: {e}");
            set_state("failed");
        }
    }
}
