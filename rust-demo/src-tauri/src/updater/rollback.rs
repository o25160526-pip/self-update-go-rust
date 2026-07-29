use std::{io, path::PathBuf, process::Command};
use tauri::{AppHandle, Manager};

pub fn rollback_to_previous(app: &AppHandle) -> Result<String, String> {
    let dir: PathBuf = app.path().app_local_data_dir().map_err(|e| e.to_string())?;
    let installer = dir.join("updates").join("previous-installer.exe");
    if !installer.exists() { return Err(format!("No cached previous installer found at {}", installer.display())); }
    Command::new(&installer).arg("/PASSIVE").spawn().map_err(|e| e.to_string())?;
    app.exit(0);
    Ok("rollback-started".into())
}

pub fn cached_installer_path(app: &AppHandle) -> Result<PathBuf, io::Error> {
    Ok(app.path().app_local_data_dir().map_err(io::Error::other)?.join("updates").join("previous-installer.exe"))
}
