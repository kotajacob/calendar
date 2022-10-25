calendar-config(5)

# NAME

calendar-config - configuration file format for *calendar*(1)

# CONFIGURATION

The configuration file is named *config.toml* and uses the easy to read TOML
format. We search for the file in your XDG config home plus "calendar", which
defaults to ~/.config/calendar on UNIX systems. A detailed configuration example
is typically included with your installation in `/usr/share/calendar`.

Each line in the config file is a pairing of a option and a value separated with
"=". Lines beginning with # are considered comments and are ignored, as are
empty lines.

# NOTE OPTIONS

*NotePath*
	NotePath should describe where the note would be located for this predefined
	time:
	January 2, 15:04:05, 2006, in time zone seven hours west of GMT
	
	Environment variable, such as $HOME may be used in the path and will be
	expanded appropriately. If the file is missing it is simply treated as an
	empty file. All other errors will return the error string itself (which is
	meant to be displayed to the user).

	Default: "$HOME/.local/share/calendar/2006-01-02.md"

*Editor*
	Used to specify an editor to use when opening notes. If unspecified, we will
	respect the environment variables VISUAL or EDITOR or fallback to using vi.

	Default: none

# COLOR OPTIONS

Colors can be specified in a few different ways. Terminal support may vary. By
default, all colors are specified using ANSI 16 which has the best support.
```
ANSI 16 colors: "5", "9", "12" (numbers 0-15)
ANSI 256 Colors: "86", "201", "202" (numbers 16-255)
True Color (16,777,216 colors; 24-bit): "#0000FF", "#04B575", "#3C3C3C"
```

*TodayColor*
	Foreground color used for the current date.

	Default: "2"

*InactiveColor*
	Foreground color used in inactive months.

	Default: "8"

# SPACING OPTIONS

We use one more padding on the left to account for the border around the preview
window. It looks a bit unbalanced without it, but feel free to tweak as you
please.

*LeftPadding*
	Overall left-side padding.

	Default: 2

*RightPadding*
	Overall right-side padding.

	Default: 1

*PreviewPadding*
	Left and right side padding for the preview window. This is the space
	between the words and the preview's border.

	Default: 1

*PreviewLeftMargin*
	Left-side margin size between the preview's border and the edge of the month
	widgets.

	Default: 3

*PreviewMinWidth*
	The minimum width of the preview window in characters. This does not include
	the width of the month widgets (which is always 20).

	Default: 40

*PreviewMaxWidth*
	The maximum width of the preview window in characters. This does not include
	the width of the month widgets (which is always 20).

	Default: 80