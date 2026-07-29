#[allow(dead_code)]
pub fn restart(app: &tauri::AppHandle) -> ! {
    app.restart()
}
