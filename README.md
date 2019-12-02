# FReD Devtools
## Setup
Please run the `./setup-git.sh` **in the `fred-devtools` (=this) folder**. It will clone the other projects into `src/`.
When working with GoLand please open this folder.

## Git Workflow

The `master` branch is protected and only approved pull requests can push to it. Most important part of
the workflow is `rebase`, heres a refresher on merging vs rebasing https://www.atlassian.com/git/tutorials/merging-vs-rebasing.

How do I push changes to the `master` branch?

1.  Switch to `master` -> `git checkout master`
2.  Update `master` -> `git pull --rebase` (ALWAYS use `rebase` when pulling!!!)
3.  Create new branch from `master` -> `git checkout -b tp/new-feature` (where 'tp/' is your own name/abbreviation)
4.  Work on branch and push changes
5.  Rebase master onto branch to not have merge conflicts later -> `git pull origin master --rebase` (AGAIN use`--rebase`)
6.  Push branch again, this time force push to include rebased master (`git push --force`)
7.  Create a pull request from gitlab.tu-berlin.de
8.  Get pull request reviewed and merge it into master

Some last words, keep pull requests small (not 100 files changed etc :D), so they are easier to review and rather create a lot of pull requests than one big