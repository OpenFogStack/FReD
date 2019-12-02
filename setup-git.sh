#!/bin/sh
cd ..
mkdir FReD
mv -rd fred-devtools FReD/
cd FReD

git clone https://gitlab.tubit.tu-berlin.de/mcc-fred/fred-example-client.git
git clone https://gitlab.tubit.tu-berlin.de/mcc-fred/fred-nase.git
git clone https://gitlab.tubit.tu-berlin.de/mcc-fred/fred-node.git
