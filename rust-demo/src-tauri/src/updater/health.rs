use rust_demo_logic::state::PersistentState;
use std::{fs, io, path::PathBuf};
use tauri::{AppHandle, Manager};

fn state_path(app: &AppHandle) -> Result<PathBuf, io::Error> {
    app.path()
        .app_local_data_dir()
        .map(|p| p.join("update-state.json"))
        .map_err(io::Error::other)
}
fn load(app: &AppHandle) -> Result<PersistentState, io::Error> {
    let p = state_path(app)?;
    match fs::read_to_string(p) {
        Ok(v) => serde_json::from_str(&v).map_err(io::Error::other),
        Err(e) if e.kind() == io::ErrorKind::NotFound => {
            Ok(PersistentState::new(crate::version::VERSION))
        }
        Err(e) => Err(e),
    }
}
fn save(app: &AppHandle, s: &PersistentState) -> Result<(), io::Error> {
    let p = state_path(app)?;
    if let Some(parent) = p.parent() {
        fs::create_dir_all(parent)?
    };
    let tmp = p.with_extension("tmp");
    fs::write(
        &tmp,
        serde_json::to_vec_pretty(s).map_err(io::Error::other)?,
    )?;
    fs::rename(tmp, p)
}
pub fn write_pending(app: &AppHandle, version: &str) -> Result<(), io::Error> {
    let mut s = load(app)?;
    s.stage(version).map_err(io::Error::other)?;
    save(app, &s)
}
pub fn complete_startup_health_check(app: &AppHandle) -> Result<(), io::Error> {
    let mut s = load(app)?;
    if s.pending.as_deref() == Some(crate::version::VERSION) {
        s.mark_healthy().map_err(io::Error::other)?
    } else {
        s.current = crate::version::VERSION.into();
        if s.last_known_good.is_empty() {
            s.last_known_good = s.current.clone()
        }
    };
    save(app, &s)
}
