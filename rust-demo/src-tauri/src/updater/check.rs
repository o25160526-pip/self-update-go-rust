// check.rs — logic kiểm tra update (sẽ implement đầy đủ ở P2-T3)
// P0: stub đơn giản để app compile được
use std::sync::OnceLock;
use tauri::AppHandle;

static STATE: OnceLock<std::sync::Mutex<String>> = OnceLock::new();

fn state_lock() -> &'static std::sync::Mutex<String> {
    STATE.get_or_init(|| std::sync::Mutex::new("checking".to_string()))
}

pub fn get_state() -> String {
    state_lock().lock().unwrap().clone()
}

pub fn set_state(s: &str) {
    *state_lock().lock().unwrap() = s.to_string();
}

/// Kiểm tra update khi app khởi động.
/// P2-T3 sẽ implement đầy đủ: query endpoint, so sánh semver, etc.
pub async fn check_on_startup(_app: AppHandle) {
    // No endpoint configured yet → up-to-date
    set_state("up-to-date");
}
