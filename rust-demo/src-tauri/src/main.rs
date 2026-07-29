#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod updater;
mod version;

use tauri::Manager;

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.iter().any(|a| a == "--version") {
        println!("{}", version::VERSION);
        return;
    }
    if args.iter().any(|a| a == "--print-update-state") {
        println!("checking");
        return;
    }
    let offline = args.iter().any(|a| a == "--offline-test");
    tauri::Builder::default()
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_process::init())
        .invoke_handler(tauri::generate_handler![
            get_version,
            get_update_state,
            rollback_previous
        ])
        .setup(move |app| {
            if let Err(e) = updater::health::complete_startup_health_check(app.handle()) {
                eprintln!("health check state failed: {e}");
            }
            let handle = app.handle().clone();
            tauri::async_runtime::spawn(async move {
                updater::check::check_on_startup(handle, offline).await
            });
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
fn get_version() -> String {
    version::VERSION.into()
}

#[tauri::command]
fn get_update_state() -> String {
    updater::check::get_state()
}

#[tauri::command]
fn rollback_previous(app: tauri::AppHandle) -> Result<String, String> {
    updater::rollback::rollback_to_previous(&app)
}
