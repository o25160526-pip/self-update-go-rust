pub struct MockResult {
    pub has_update: bool,
}
pub fn check() -> MockResult {
    MockResult { has_update: true }
}
