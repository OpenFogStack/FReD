#!/bin/sh
# Execute this in the fred-devtools folder!
cd ..
echo "Moving folder fred-devtools to FReD"
mv fred-devtools FReD
cd FReD

mkdir src
cd src
echo "Cloning into fred-example, fred-nase and fred-node"
git clone git@gitlab.tubit.tu-berlin.de:mcc-fred/fred-example-client.git
git clone git@gitlab.tubit.tu-berlin.de:mcc-fred/fred-nase.git
git clone git@gitlab.tubit.tu-berlin.de:mcc-fred/fred-node.git

echo "You can now open this project (folder FReD) in GoLand"
