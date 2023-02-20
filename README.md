# FReD: Fog Replicated Data

[![pipeline status](https://git.tu-berlin.de/mcc-fred/fred/badges/main/pipeline.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![coverage report](https://git.tu-berlin.de/mcc-fred/fred/badges/main/coverage.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/git.tu-berlin.de/mcc-fred/fred)](https://goreportcard.com/report/git.tu-berlin.de/mcc-fred/fred)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/git.tu-berlin.de/mcc-fred/fred)

**FReD** is a distributed middleware for **F**og **Re**plicated **D**ata.
It abstracts data management for fog-based applications by grouping data into _keygroups_, each keygroup a set of key-value pairs that can be managed independently.
Applications have full control over keygroup replication: replicate your data where you need it.

FReD is maintained by [Tobias Pfandzelter, Trever Schirmer, and Nils Japke of the Mobile Cloud Computing research group at Technische Universität Berlin and Einstein Center Digital Future in the scope of the FogStore project](https://www.tu.berlin/en/mcc).
Funded by the Deutsche Forschungsgemeinschaft (DFG, German Research Foundation) -- 415899119.

FReD is open-source software and contributions are welcome.
You can find the complete documentation on the [FReD website](https://openfogstack.github.io/fred).
All contributions should be submitted as merge requests on the [main repository on the TU Berlin GitLab](https://git.tu-berlin.de/mcc-fred/fred) and are subject to review by the maintainers.
Check out the [Contributing](#contributing) section for more information.
