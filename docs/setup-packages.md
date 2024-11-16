# Setup of Packages

## Oh My Zsh

> Command line tools.

https://ohmyz.sh/#install

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

And use these settings

```shell
# Add to ~/.zshrc
source ~/Developer/jterrazz-configuration/scripts/main.sh
```

## Brew

> Package manager for MacOS.

https://brew.sh

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Nvm

> Node version manager.

https://github.com/nvm-sh/nvm

```shell
brew install nvm
```

## Git

### Create the SSH key

```shell
ssh-keygen -t ed25519 -C "j2aa@proton.me"
eval "$(ssh-agent -s)"
```

Link this key to the GitHub host.

```shell
# vim ~/.ssh/config

Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
```

### Use the SSH key

```shell
ssh-add --apple-use-keychain ~/.ssh/id_ed25519
```