---
layout: default
title: Contributing
nav_order: 2
---

## Contributing

For development, it is recommended to install [GoLand](https://www.jetbrains.com/go/).

### Git Workflow

Setup git environment with `sh ./ci/env-setup.sh` (installs git hooks). Be sure to have Go (>1.16) installed.

The `main` branch is protected and only approved pull requests can push to it.
Most important part of the workflow is `rebase`, [here's](https://www.atlassian.com/git/tutorials/merging-vs-rebasing) a refresher on merging vs rebasing.

1. Switch to `main` → `git checkout main`
2. Update `main` → `git pull --rebase` (ALWAYS use `rebase` when pulling!)
3. Create new branch from `main` → `git checkout -b tp/new-feature` (where `tp` is your own name/abbreviation)
4. Work on branch and push changes
5. Rebase `main` onto branch to not have merge conflicts later → `git pull origin main --rebase` (AGAIN use`--rebase`)
6. Push branch again, this time force push to include rebased `main` (`git push --force`)
7. Create a pull request from `git.tu-berlin.de`
8. Get pull request reviewed and merge it into `main`

Some last words, keep pull requests small (not 100 files changed...), so they are easier to review and rather create a lot of small pull requests than one big.

### Code Quality and Testing

In order to keep our code clean and working, we provide a number of test suites and support a number of code quality tools.

#### Static Analysis

Static analysis analyses the code in the repository without actually executing any of it.

##### Compiling

Before anything else, code must of course compile to be considered valid.
To compile the main code, use `make`, which builds the main `fred` software in the `./cmd/frednode` folder.
It first fetches dependencies and then builds the software.
Any compiler warnings and errors produced by the build process make this code invalid.

##### Linting

We use the [`golangci-lint`](https://github.com/golangci/golangci-lint) to run a number of linting tasks on our code.
Check [the documentation](https://golangci-lint.run/usage/install/#local-installation) to install it on your machine.

Once the utility is available to you, run `make lint` to lint all Go code files in the repository.
This uses the [default list of linters](https://golangci-lint.run/usage/linters/#enabled-by-default-linters), finds some basic errors and style faults, and marks them.
These linting errors mostly don't make the code invalid but ignoring them can lead to bad code quality.

##### Mega-Linting

Additionally, we also provide the `make megalint` command in our Makefile.
This runs a number of additional checks on the code, yet passing these checks is not mandatory for code to be merged into the repository.

#### Unit Tests

For a number of packages we also include unit tests to test some functionality.
This allows testing basic functionality, test for race conditions, and memory sanitation.
Use `make test`, `make race`, and `make msan` to execute all tests.

To get a coverage report, you can use `make coverage` and `make coverhtml`, depending on your preference.

#### System Tests

There are two system tests to test functionality of a FReD deployment as a whole.
This is part of a TDD approach where tests can be defined first, and the software is refined until it completes all tests.
All system tests can be found in `./tests`.
All tests require Docker and Docker Compose to work.

##### 3 Node Test

The "3 node test" starts a FReD deployment of three FReD nodes and runs a client against the FReD cluster that validates different functionalities.
It uses Docker compose and can thus easily be started with `make 3n-all`.

The deployment comprises a single `etcd` Docker container as a NaSe, a simple trigger node, two FReD nodes that each comprise only a single machine (node _B_ and _C_) with a storage server, and a distributed FReD node _A_ that comprises three individual FReD machines behind a `fredproxy` sharing a single DynamoDB storage server.
All machines are connected over a Docker network.

The test client runs a number of operations against the FReD deployment and outputs a list of errors.
The complete code for the test client can be found in `./tests/3NodeTest`.

When the debug log output of the individual nodes is not enough to debug an issue, it is also possible to connect a `dlv` debugger directly to FReD node _B_ to set breakpoints or step through code.
This is currently configured to use the included debugger in the GoLand IDE.
Further information can be found in the 3 node test documentation.

##### Failing Node Test

As FReD is a distributed system, it is important to also test the impact of a failing node in the deployment.
The "Failing Node Test" allows this.
It starts the same deployment as in the 3 node test and runs a number of queries before killing one of the nodes and starting it back up.
It uses the Docker API to destroy and start the corresponding containers.

The code can be found in `./tests/FailingNodeTest` but can be started with `make failtest` in `./tests/3NodeTest/` after a deployment has been created with `make fred`.

##### ALExANDRA Test

The ALExANDRA test tests a limited amount of middleware functionality.
Use `make alexandratest` to run it.
The complete code can be found and extended in `./tests/AlexandraTest`.

##### Consistency Test

The consistency test tests consistency guarantees provided by the middleware.
In `./tests/consistency` run `bash ./run-cluster.sh [NUM_NODES] [NUM_CLIENTS]` and specify the number of FReD nodes and clients.

#### Cluster

You can easily set up a cluster of FReD nodes by using the `run-cluster.sh` script in the `cluster/` folder.
Simply run `bash run-cluster.sh [NUM_NODES]` to spawn up to 263 FReD nodes.

#### Profiling

FReD supports CPU and memory profiling for the main `frednode` binary.
Use the `--cpuprofile` and `--memprofile` flags in addition to your other flags to enable profiling.
Keep in mind that this may have an impact on performance in some cases.

```sh
# start fred with profiling
$ ./frednode --cpuprofile fredcpu.pprof --memprof fredmem.pprof [ALL_YOUR_OTHER_FLAGS...]

# run tests, benchmarks, your application, etc.
# then quit fred with SIGINT or SIGKILL and your files will be written
# open pprof files and convert to pdf (note that you need graphviz installed)
# you also need to provide the path to your frednode binary
$ go tool pprof --pdf ./frednode fredcpu.pprof > cpu.pdf
$ go tool pprof --pdf ./frednode fredmem.pprof > mem.pdf

```
