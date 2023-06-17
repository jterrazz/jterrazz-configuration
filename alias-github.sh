alias.lola=log --graph --pretty='format:%C(auto)%h %d %s %C(green)%an%C(bold blue) %ad' --all --date=relative
alias.lol=log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
alias.oops=!f(){ if [ "$1" == '' ]; then   git commit --amend --no-edit; else   git commit --amend "$@"; fi;}; f
alias.fop=fetch origin --prune
alias.pp=pull --prune
alias.pfnov=push -f --no-verify
alias.pnov=push --no-verify
alias.wip-branch=for-each-ref --sort='authordate:iso8601' --format=' %(color:green)%(authordate:relative)%09%(color:white)%(refname:short)' refs/heads
alias.co=checkout
alias.cob=checkout -b
alias.br=branch
alias.ci=commit
alias.st=status
alias.mas=rebase main
alias.unstage=reset HEAD --
alias.sp=switch -
alias.fp=!git fetch -p; git pull
alias.wip=!git add --all; git ci -m WIP
alias.unwip=!git reset --soft HEAD~1; git unstage
alias.prunel=!git branch --merged | grep -v "\*" | grep -v main | xargs -n 1 git branch -d
alias.ri=!sh -c 'git rebase -i HEAD~$(git rev-list --count $(git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev)..HEAD)'
alias.changelog=!f() { git log --topo-order --no-merges --pretty=format:"%s (%an)" $1..HEAD; }; f
alias.feat=!f(){ git add . && git commit -m "[FEATURE] $1"; };f
alias.chore=!f(){ git add . && git commit -m "[CHORE] $1"; };f
alias.fix=!f(){ git add . && git commit -m "[BUGFIX] $1"; };f
