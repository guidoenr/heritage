#!/bin/bash

cd ~/Documents/heritage

# if no changes
if [[ "$(git pull)" == *"Already up to date."* ]]; then
    # RUN
    make run
else
    make build
    make run
fi

sleep 1

