use std::{fs, io, path::PathBuf};

use tauri::{AppHandle, Manager};

pub fn rollback_to_previous(app: &AppHandle) -> Result<String, String> {
    let path = cached_installer_path(app).map_err(|e| e.to_string())?;
    if !path.exists() {
        return Err(format!("cached installer not found: {}", path.display()));
    }
    std::process::Command::new(&path)
        .arg("/S")
        .status()
        .map_err(|e| e.to_string())?;
    Ok(format!("rollback installer started: {}", path.display()))
}

#[allow(dead_code)]
pub fn cached_installer_path(app: &AppHandle) -> Result<PathBuf, io::Error> {
    let dir = app.path().app_data_dir()?;
    Ok(dir.join("previous").join("rust-demo-previous-setup.exe"))
}
