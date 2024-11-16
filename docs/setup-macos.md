# Setup of MacOS

## Dock

- [x] Spacing for icons

```shell
defaults write com.apple.dock persistent-apps -array-add '{"tile-type"="small-spacer-tile";}'; killall Dock
```
