# GitHub Security

## Create the SSH key

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

## Use the SSH key
```shell
ssh-add --apple-use-keychain ~/.ssh/id_ed25519
```