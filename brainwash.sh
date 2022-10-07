#!/bin/bash

USER="$1"

source /home/$USER/.profile

# You may need to edit this to point to your working directory of the repo
cd /home/$USER/mister-brainwash

make run
