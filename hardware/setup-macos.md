# MacOS Setup

## Dock

```shell
defaults write com.apple.dock persistent-apps -array-add '{"tile-type"="small-spacer-tile";}'; killall Dock
```