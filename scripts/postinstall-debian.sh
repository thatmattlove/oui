#!/bin/bash

find_path() {
    local paths=($(dpkg -L oui))
    local path=${paths[-1]}
    echo $path
    return 0
}

BIN=$(find_path)

$BIN update
exit 0
