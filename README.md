# Punos - DJing tool for CUI users

As is well known, people sometimes want to DJ, and Engineers always use CUI are no exception.
Punos is tool for DJing on CUI target such people and written in Go.

![](https://i.imgur.com/GLg4Nae.gif)

## Features
- Play music
  - Play/Pause
  - FastForward/Rewind
  - Volume Control
  - Speed Control
  - Show position
- Cue
- Show audio wave form

## Support Audio Format
- mp3

## Support OS
- mac

## Installation
```
$ go get github.com/amaretto/punos/cmd/punos
```

## Usage
1. create directory named "mp3" and put music files to it
2. run the command in same path with "mp3" directory(At first time, It may take some time to generate waveform data)
```
$ punos
```

If you want to play next music, you can use other windows

![](https://i.imgur.com/OejXOfI.png)

## Keybindings
### PunosPanel(Main Screen)
| key         | description          |
|-------------|----------------------|
| Space       | Play/Pause           |
| ESC         | Quit                 |
| w           | Fast Forward         |
| q           | Rewind               |
| s           | Volume Up            |
| a           | Volume Down          |
| x           | Speed Up             |
| z           | Speed Down           |
| f           | Switch to Load Panel |

### LoadPanel
| key         | description          |
|-------------|----------------------|
| ESC         | Quit                 |
| j           | Down Cursor          |
| k           | Up Cursor            |
| l           | Load Music           |
| a           | Analyze All Music    |

## License
MIT

## Author
[amaretto](https://github.com/amaretto)
