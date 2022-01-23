#!/bin/bash

BIN="$(which oui)"

if [[ "$?" == "1" ]]; then
    source ~/.profile
    BIN="$(which oui)"
    if [[ "$?" == "1" ]]; then
        echo $'Unable to locate oui in $PATH. You\'ll need to run \'oui update\' once oui has been added to $PATH'
        exit 0
    fi
fi

$BIN update
exit 0
