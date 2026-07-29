// Tauri updater không có rollback built-in. State adapter giữ last-known-good;
// production cần cache installer cũ và launcher ngoài tiến trình để thực thi rollback.
pub use rust_demo_logic::state::PersistentState;
