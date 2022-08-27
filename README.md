# `calendar`

A TUI version of the classic `cal` program.
```
     August 2022    
Su Mo Tu We Th Fr Sa
    1  2  3  4  5  6
 7  8  9 10 11 12 13
14 15 16 17 18 19 20
21 22 23 24 25 26 27
28 29 30 31         
```

## Navigation
Using the arrow keys, hjkl, or the mouse you can navigate around the calendar.

## Display
The program expands to use as much terminal space as you provide. More vertical
space allows `calendar` to show the previous and next months stacked on top of
each other in a lighter grey color (configurable). The current day is shown in
green (also configurable).

Next to the month widget(s) you will see a preview file for that day's note
file. You can configure a path to store these notes and press `enter` to open
them in your favorite editor.
```
       July 2022           Jazzi and I had breakfast at adjo. It was
  Su Mo Tu We Th Fr Sa     rainy and quite cozy and we managed to get
                  1  2     there early enough that it wasn't busy
   3  4  5  6  7  8  9     yet. When we got home I spent pretty much
  10 11 12 13 14 15 16     the whole day working on getting my
  17 18 19 20 21 22 23     raspberry pi setup on our tv... It didn't
  24 25 26 27 28 29 30     go very well. I ran into loads of issues.
  31                       It was already running alpine linux which
      August 2022          is nice, but I still don't really know a
  Su Mo Tu We Th Fr Sa     good way to list all your manually
      1  2  3  4  5  6     installed packages and it had gotten quite
   7  8  9 10 11 12 13     messy from me using it for various
  14 15 16 17 18 19 20     experiments.
  21 22 23 24 25 26 27
  28 29 30 31              For whatever reason I just could not get
                           kodi either from xorg or wayland or even
     September 2022        the direct rendering "gbm" method. I
  Su Mo Tu We Th Fr Sa     didn't have a keyboard so I'm trying to do
               1  2  3     this all over ssh which means it's really
   4  5  6  7  8  9 10     annoying to get any sort of error messages
  11 12 13 14 15 16 17     and I had to constantly keep rebooting the
  18 19 20 21 22 23 24     thing. I should've just gotten a keyboard
  25 26 27 28 29 30        before getting started. Eventually I
```

## Configuration
I've written this tool to be fairly configurable. Padding, colors, path to your
notes and of course your editor can all be configured in
`XDG_CONFIG_HOME/calendar/config.toml`. Normally, on unix systems that would be
`$HOME/.config/config.toml`. This repository contains an example config.toml
showing every option with detailed explanatory comments.

## Development
### Hot Reload
Sadly, `entr` doesn't seem to work with tui programs. So, I've been needing to
do this insane hack of having it spawn a terminal window running `go run .`
which works better than you might expect, but is still kind of annoying:
```
autostash alias rp='fd -e go | entr -r alacritty --class "Alacritty-entr,Alacritty-entr" -o window.position.x=1380 -o window.position.y=82 -e go run .'
```

# Author
Written and maintained by Dakota Walsh.
Up-to-date sources can be found at https://git.sr.ht/~kota/calendar/

# License
GNU GPL version 3 only, see LICENSE.
