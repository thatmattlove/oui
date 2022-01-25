#!/bin/bash

BIN="$(which oui)"

if command -v oui >/dev/null; then
    oui update
    exit 0
fi

echo $'Unable to run `oui update` automatically'
exit 1
