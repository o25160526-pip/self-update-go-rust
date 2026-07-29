# Test Cases

> 7 nhóm test A-G theo prompt §12.
> Chi tiết đầy đủ sẽ hoàn thiện ở `docs/testing.md` (P5-T3).
> Mỗi test phải có: mã, setup, command, expected result, log cần kiểm tra, cleanup.

## Nhóm A — App cơ bản
- A01: clean build Go
- A02: clean build Rust
- A03: app hiển thị đúng version
- A04: `--version` trả đúng version
- A05: app khởi động offline
- A06: app không ghi secret vào log

## Nhóm B — Update bình thường
- B01: `1.0.0`, không có release mới, app báo up-to-date
- B02: `1.0.0`, có `1.0.1`, app báo update available
- B03: auto download không cần click
- B04: hash đúng, signature đúng, install thành công
- B05: app tự restart
- B06: sau restart hiển thị `1.0.1`
- B07: state ghi `last-known-good=1.0.1`

## Nhóm C — Lỗi mạng
- C01: DNS fail
- C02: timeout
- C03: HTTP 404
- C04: HTTP 500
- C05: download bị ngắt giữa chừng
- C06: resume hoặc retry đúng policy
- C07: app vẫn chạy version cũ, không bị hỏng

## Nhóm D — Bảo mật
- D01: SHA-256 sai, update bị từ chối
- D02: signature sai, update bị từ chối
- D03: manifest bị sửa, update bị từ chối
- D04: version thấp hơn hiện tại, bị từ chối
- D05: artifact không đúng OS/architecture, bị từ chối
- D06: private key không xuất hiện trong artifact/log
- D07: client không cần GitHub token với public release

## Nhóm E — Restart & rollback
- E01: app đang giữ file, helper chờ app thoát
- E02: app crash trước khi health-check
- E03: app mới treo quá timeout
- E04: local service không khởi động
- E05: rollback về last-known-good
- E06: không rollback lặp vô hạn
- E07: máy tắt giữa lúc install, khởi động lại vẫn phục hồi được
- E08: rollback không xóa nhầm user data/config

## Nhóm F — Quyền Windows
- F01: cài per-user, update không cần admin
- F02: cài machine-wide, UAC chỉ xuất hiện ở bước cần thiết
- F03: user không có quyền ghi thư mục, lỗi rõ ràng
- F04: artifact được tải vào thư mục temp có permission phù hợp
- F05: code signing/SmartScreen behavior được ghi nhận

## Nhóm G — Phát hành
- G01: tag `v1.0.0` tạo release
- G02: tag `v1.0.1` tạo artifact mới
- G03: artifact Go và Rust không ghi đè nhau
- G04: checksum upload đúng
- G05: manifest trỏ đúng artifact
- G06: workflow fail khi thiếu signing secret
- G07: release có release notes
