#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod version;
mod updater;

use tauri::Manager;

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_process::init())
        .invoke_handler(tauri::generate_handler![get_version, get_update_state])
        .setup(|app| {
            let handle = app.handle().clone();
            tauri::async_runtime::spawn(async move {
                updater::check::check_on_startup(handle).await;
            });
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
fn get_version() -> String {
    version::VERSION.to_string()
}

#[tauri::command]
fn get_update_state() -> String {
    updater::check::get_state()
}
