#!/bin/bash

find_path() {
    local paths=($(rpm -ql oui | grep '/oui$'))
    local path=${paths[0]}
    return path
}

BIN=$(find_path)

$BIN update
exit 0
