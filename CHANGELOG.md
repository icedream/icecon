# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added
- Up and down arrow keys can now be used to scroll through command history to reissue already typed commands. (Windows UI) ([#10](https://github.com/icedream/icecon/pull/10), thanks to @TheIndra55)
- Current server address in title bar. (Windows UI)

### Changed
- Binaries are now compiled statically and using Go 1.8.
- Reuse server address as typed in by the user in connect dialog instead of using resolved IP address. (Windows UI)
- Update copyright text.

## [1.0.0] - 2016-05-07
### Added
- Add fully working command line flags, see the help text that can be called by running IceCon with `--help` in a console
- Add graphical UI for Windows (`--gui` or automatically shown when run without parameters)
- Add netcat-style console interface (can be used like netcat to pipe through commands)
- Add script-friendly command line interface (`icecon -c <your command here> <server:port> <password>`)

[Unreleased]: https://github.com/icedream/icecon/compare/v1.0.0...develop
[1.0.0]: https://github.com/icedream/icecon/releases/tag/v1.0.0
