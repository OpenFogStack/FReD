# FReD: Fog Replicated Data

[![pipeline status](https://git.tu-berlin.de/mcc-fred/fred/badges/main/pipeline.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![coverage report](https://git.tu-berlin.de/mcc-fred/fred/badges/main/coverage.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/git.tu-berlin.de/mcc-fred/fred)](https://goreportcard.com/report/git.tu-berlin.de/mcc-fred/fred)
[![Godoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/git.tu-berlin.de/mcc-fred/fred)
[![Docs](https://img.shields.io/badge/-documentation-informational.svg)](https://openfogstack.github.io/FReD/)

**FReD** is a distributed middleware for **F**og **Re**plicated **D**ata.
It abstracts data management for fog-based applications by grouping data into _keygroups_, each keygroup a set of key-value pairs that can be managed independently.
Applications have full control over keygroup replication: replicate your data where you need it.

FReD is maintained by [Tobias Pfandzelter, Trever Schirmer, and Nils Japke of the Mobile Cloud Computing research group at Technische Universität Berlin and Einstein Center Digital Future in the scope of the FogStore project](https://www.tu.berlin/en/mcc).
Funded by the Deutsche Forschungsgemeinschaft (DFG, German Research Foundation) -- 415899119.

If you use this software in a publication, please cite it as:

### Text

T. Pfandzelter, N. Japke, T. Schirmer, J. Hasenburg, and D. Bermbach, **Managing Data Replication and Distribution in the Fog with FReD**, Software: Practice and Experience, Jul. 2023.

### BibTeX

```bibtex
@article{pfandzelter2023fred,
    author = "Pfandzelter, Tobias and Japke, Nils and Schirmer, Trever and Hasenburg, Jonathan and Bermbach, David",
    title = "Managing Data Replication and Distribution in the Fog with FReD",
    journal = "Software: Practice and Experience",
    month = jul,
    year = 2023,
    issn = "0038-0644",
    url = "https://doi.org/10.1002/spe.3237",
    doi = "10.1002/spe.3237",
    publisher = "Wiley",
    address = "Hoboken, NJ, USA",
}
```

For a full list of publications, please see [our website](https://www.tu.berlin/en/mcc/research/publications).

### License

The code in this repository is licensed under the terms of the [MIT](./LICENSE) license.
FReD is open-source software and contributions are welcome.
You can find the complete documentation on the [FReD website](https://openfogstack.github.io/FReD).
All contributions should be submitted as merge requests on the [main repository on the TU Berlin GitLab](https://git.tu-berlin.de/mcc-fred/fred) and are subject to review by the maintainers.
Check out the [Contributing](#contributing) section for more information.
