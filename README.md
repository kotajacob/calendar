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
```
       July 2022     
 Su Mo Tu We Th Fr Sa
                 1  2
  3  4  5  6  7  8  9
 10 11 12 13 14 15 16
 17 18 19 20 21 22 23
 24 25 26 27 28 29 30
 31                  
      August 2022    
 Su Mo Tu We Th Fr Sa
     1  2  3  4  5  6
  7  8  9 10 11 12 13
 14 15 16 17 18 19 20
 21 22 23 24 25 26 27
 28 29 30 31         
                     
    September 2022   
 Su Mo Tu We Th Fr Sa
              1  2  3
  4  5  6  7  8  9 10
 11 12 13 14 15 16 17
 18 19 20 21 22 23 24
 25 26 27 28 29 30
```

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
GNU GPL version 3 or later, see LICENSE.
