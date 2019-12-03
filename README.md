# FReD 

[![pipeline status](https://gitlab.tubit.tu-berlin.de/mcc-fred/fred/badges/master/pipeline.svg)](https://gitlab.tubit.tu-berlin.de/mcc-fred/fred/commits/master)
[![coverage report](https://gitlab.tubit.tu-berlin.de/mcc-fred/fred/badges/master/coverage.svg)](https://gitlab.tubit.tu-berlin.de/mcc-fred/fred/commits/master)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

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