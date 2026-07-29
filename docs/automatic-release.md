# Automatic release on every main push

Every push to `main` now creates the same deterministic patch release in both Go and Rust workflows:

`1.0.<git commit count>`

Both workflows derive the version from the checked-out commit, so they race safely on one tag and upload separate assets to it. Release branches and explicit `v*` tags preserve their requested version.

Pull requests do not publish releases. They run smoke/build validation only.

Do not rewrite a published release. Fix forward with the next commit and let the next patch release carry the fix.
