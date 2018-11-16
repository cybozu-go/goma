# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.2] - 2018-11-16
- Handle renaming of cybozu-go/cmd to [cybozu-go/well][well]
- Introduce support for Go modules

## [1.0.1] - 2016-08-24
### Changed
- Fix for cybozu-go/cmd v1.1.0.

## [1.0.0] - 2016-08-22
### Added
- goma now adopts [github.com/cybozu-go/cmd][cmd] framework.  
  As a result, commands implement [the common spec][spec].
- [actions/exec] new parameter "debug" to log outputs on failure.

[well]: https://github.com/cybozu-go/well
[cmd]: https://github.com/cybozu-go/cmd
[spec]: https://github.com/cybozu-go/cmd/blob/master/README.md#specifications
[Unreleased]: https://github.com/cybozu-go/goma/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/cybozu-go/goma/compare/v1.0.0...v1.0.1
[1.0.1]: https://github.com/cybozu-go/goma/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/cybozu-go/goma/compare/v0.1...v1.0.0
