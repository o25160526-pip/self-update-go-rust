use std::fs;
use std::path::Path;

/// Kich thuoc icon duoc sinh tu dong (px).
const ICON_SIZE: u32 = 64;

/// Sinh `icons/icon.ico` (64x64, 32bpp BMP) neu file chua ton tai.
///
/// `tauri-build`, `tauri::generate_context!` va NSIS bundler deu yeu cau mot file
/// `.ico` hop le khi build tren Windows. Repo chi commit file text nen icon duoc
/// sinh ngay tai thoi diem build, giup build xanh tren CI lan may Windows 11.
fn ensure_windows_icon(path: &Path) -> std::io::Result<()> {
    if path.exists() {
        return Ok(());
    }
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)?;
    }

    // Pixel data: BGRA, luu bottom-up theo chuan BMP.
    let mut pixels: Vec<u8> = Vec::new();
    for y in 0..ICON_SIZE {
        let row = ICON_SIZE - 1 - y;
        for x in 0..ICON_SIZE {
            pixels.push(235);
            pixels.push((60 + row * 2) as u8);
            pixels.push((20 + x) as u8);
            pixels.push(0xFF);
        }
    }

    // AND mask 1bpp, moi dong padding len boi so cua 4 byte.
    let mask_stride = (((ICON_SIZE + 31) / 32) * 4) as usize;
    let mask = vec![0u8; mask_stride * ICON_SIZE as usize];

    let image_bytes = (pixels.len() + mask.len()) as u32;
    let mut dib: Vec<u8> = Vec::new();
    dib.extend_from_slice(&40u32.to_le_bytes()); // biSize
    dib.extend_from_slice(&(ICON_SIZE as i32).to_le_bytes()); // biWidth
    dib.extend_from_slice(&(ICON_SIZE as i32 * 2).to_le_bytes()); // biHeight + mask
    dib.extend_from_slice(&1u16.to_le_bytes()); // biPlanes
    dib.extend_from_slice(&32u16.to_le_bytes()); // biBitCount
    dib.extend_from_slice(&0u32.to_le_bytes()); // biCompression = BI_RGB
    dib.extend_from_slice(&image_bytes.to_le_bytes()); // biSizeImage
    dib.extend_from_slice(&0u32.to_le_bytes()); // biXPelsPerMeter
    dib.extend_from_slice(&0u32.to_le_bytes()); // biYPelsPerMeter
    dib.extend_from_slice(&0u32.to_le_bytes()); // biClrUsed
    dib.extend_from_slice(&0u32.to_le_bytes()); // biClrImportant
    dib.extend_from_slice(&pixels);
    dib.extend_from_slice(&mask);

    let mut ico: Vec<u8> = Vec::new();
    ico.extend_from_slice(&0u16.to_le_bytes()); // reserved
    ico.extend_from_slice(&1u16.to_le_bytes()); // type = icon
    ico.extend_from_slice(&1u16.to_le_bytes()); // so luong anh
    ico.push(ICON_SIZE as u8); // width
    ico.push(ICON_SIZE as u8); // height
    ico.push(0); // palette
    ico.push(0); // reserved
    ico.extend_from_slice(&1u16.to_le_bytes()); // planes
    ico.extend_from_slice(&32u16.to_le_bytes()); // bits per pixel
    ico.extend_from_slice(&(dib.len() as u32).to_le_bytes());
    ico.extend_from_slice(&22u32.to_le_bytes()); // offset toi anh dau tien
    ico.extend_from_slice(&dib);

    fs::write(path, ico)
}

fn main() {
    let icon = Path::new("icons").join("icon.ico");
    if let Err(e) = ensure_windows_icon(&icon) {
        panic!("khong sinh duoc {}: {e}", icon.display());
    }
    tauri_build::build()
}
