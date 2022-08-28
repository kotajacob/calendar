calendar(1)

# NAME

calendar - kota's little TUI calendar program

# SYNOPSIS

_calendar_

A TUI version of the classic *cal*(1) program with the ability to create, edit,
and view note files for each day. It can be used to keep a daily journal, plan
out future events, or to simply browse an interactive calendar.

# CONTROLS

Using the arrow keys, hjkl, or the mouse you can navigate around the calendar.
Pressing tab will switch your focus to the preview window so you scroll it up
and down. Pressing enter will open the selected day's note in your editor.

# DISPLAY

The program expands to use as much terminal space as you provide. More vertical
space allows *calendar* to show the previous and next months stacked on top of
each other in a lighter grey color (configurable). The current day is shown in
green (also configurable).

If your terminal is wide enough you will see a preview window next to the month
widget(s) with the selected day's note. You can configure a path to store these
notes and press enter to open them in your favorite editor. See
*calendar-config*(5) for configuration details.

# AUTHORS

Maintained by Dakota Walsh <kota at nilsu.org>. Up-to-date sources can be found
at https://git.sr.ht/~kota/calendar/
