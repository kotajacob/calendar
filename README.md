# `calendar`

I've gotten started by simply cloning the look of the `cal` program, but in a
tui.
```
     August 2022    
Su Mo Tu We Th Fr Sa
    1  2  3  4  5  6
 7  8  9 10 11 12 13
14 15 16 17 18 19 20
21 22 23 24 25 26 27
28 29 30 31         
```

## Journal Feature
My idea is that when you click or press enter on a particular day it will open
your text editor (similar to how aerc works) from right inside this tool. The
exact file path that opens will be configurable. For example I'd like mine to
open one of the date entries in my memex such as 2022-08-04 which is stored in
`~/docs/memex/journal/2022-08-04.md`. I can use go's templating and date
formatting features for that.

## Display
My other thought is that this program should a bit more of the space if you use
a large terminal size and center itself nicely. If you have more vertical space
it will show the previous and next months stacked on top of each other, but in
the lighter grey color. Additionally, if you have more horizontal space it could
show a preview of the selected date's entry off to the right side. Perhaps
something like this:
```
       July 2022         
 Su Mo Tu We Th Fr Sa    
                 1  2    Fortunately, I was only about 50 minutes late. We began  
  3  4  5  6  7  8  9    in an intense battle with one of those spooky shadowy    
 10 11 12 13 14 15 16    cats that captured finn last session was fighting with    
 17 18 19 20 21 22 23    the rest of the group when I came through the portal.
 24 25 26 27 28 29 30                                                  
 31                                                                               
      August 2022        During the battle the cat seemed to summon these creepy  
 Su Mo Tu We Th Fr Sa    tall grey men that charged at us, it was quite an        
     1  2  3 [4] 5  6    intense fight and although we survived we're beaten      
  7  8  9 10 11 12 13    bloody and will need a few days of rest to heal up.      
 14 15 16 17 18 19 20                                                             
 21 22 23 24 25 26 27    Despite our injuries finn encouraged the rest to quickly 
 28 29 30 31             look at the magic circle down beyond the alter's portal. 
                         I was lowered on a rope by moxiegordon, but we    
    September 2022       expected the room to be upsidedown as it was. For        
 Su Mo Tu We Th Fr Sa    whatever reason it was no longer upsidedown and I        
              1  2  3    dropped from the ceiling and nearly slammed into the     
  4  5  6  7  8  9 10    floor. It was very unexpected and really knocked the     
 11 12 13 14 15 16 17    wind out of me.                                          
 18 19 20 21 22 23 24 
 25 26 27 28 29 30
```

For preview file, I think if the given filename paths end in `.md` we could do
paragraph rewrapping by default. Otherwise we will do the somewhat ugly hard
wrapping. There can be an option of course to manually specify it or not, but as
a default I think this is pretty smart.

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
