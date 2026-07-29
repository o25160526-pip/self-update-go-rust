use tauri::AppHandle;
use tauri_plugin_updater::Update;

pub async fn download_install_restart(app: AppHandle, update: Update) -> Result<(), String> {
    super::check::set_state("downloading");
    let version = update.version.clone();
    update
        .download_and_install(|_, _| {}, || super::check::set_state("verifying"))
        .await
        .map_err(|e| e.to_string())?;
    super::health::write_pending(&app, &version).map_err(|e| e.to_string())?;
    super::check::set_state("restarting");
    app.restart();
}
