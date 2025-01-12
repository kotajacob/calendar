calendar(1) [""["calendar-VERSION"]]

# NAME

calendar - kota's little TUI calendar program

# SYNOPSIS

*calendar* [[[_day_] _month_] _year_]

*calendar* [_timestamp_|_monthname_]

A TUI version of the classic *cal*(1) program with the ability to create, edit,
and view note files for each day. It can be used to keep a daily journal, plan
out future events, or to simply browse an interactive calendar. If no date is
given, the current time is selected.

If giving a timestamp it should be in the form YYYY-MM-DD or DD MM YYYY.

# CONTROLS

The default controls are below. See *calendar-config*(5) for configuration
details.

[[ *Select*
:< hjkl, arrow keys, or mouse
|  *Focus preview*
:< tab
|  *Enable/disable preview*
:< p
|  *Scroll preview*
:< jk, up/down (if focused)
|  *Edit note*
:< enter
|  *Copy date*
:< y
|  *Goto last Sunday*
:< b, H
|  *Goto next Sunday*
:< w
|  *Goto next Saturday*
:< e, L
|  *Go up one month*
:< ctrl+u
|  *Go down one month*
:< ctrl+d

# DISPLAY

The program expands to use as much terminal space as you provide. More vertical
space allows *calendar* to show the previous and next months stacked on top of
each other in a lighter grey color (configurable). The current day is shown in
green (also configurable).

If your terminal is wide enough you will see a preview window next to the month
widget(s) with the selected day's note. You can configure a path to store these
notes and press enter to open them in your favorite editor. See
*calendar-config*(5) for configuration details.

If you press p (configurable) to disable the preview and your terminal is wide
enough you will be shown a full year view.

# HOLIDAYS

You can configure a list of yearly dates, such as birthdays, holidays, or other
important re-occuring events which will be displayed in a configurable color
with an optional message. See *calendar-config*(5) for configuration details.

# SEE ALSO

*calendar-config*(5)

# AUTHORS

Maintained by Dakota Walsh <kota at nilsu.org>. Up-to-date sources can be found
at https://git.sr.ht/~kota/calendar/
