#!/bin/bash

find_path() {
    local paths=($(dpkg -L oui))
    local path=${paths[-1]}
    return path
}

BIN=$(find_path)

$BIN update
exit 0
