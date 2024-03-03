# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### Added
- Add keyword display feature.
- Add holiday display feature.
- Ability to style days which have notes written differently.

### Changed
- Changed color configuration settings to general style configuration.
- Left and Right in year view can now move across to the next month.

## [0.2.0]
### Added
- Full year view mode.
- Ability to configure controls.
- Mousewheel support.
- A help menu: "?".
- Toggle preview visibility: "p".
- Copy selected date to clipboard: "y".
- Back to start of week: "b" / "H".
- End of week: "e" / "L".
- Start of next week: "w".
- Jump down one month: "ctrl+d".
- Jump up one month: "ctrl+u".
- Selecting a date with cli args.

## [0.1.0]
### Added
- Display one/three months at a time depending on terminal height.
- Display a note file for selected day depending on terminal width.
- Select different day with hjkl, arrow keys or the mouse.
- Edit selected note file in editor.
- Update "today" every night when the new day rolls over.
- TOML config file.
- Man pages detailing usage and configuration.
