#!/bin/bash

find_path() {
    local paths=($(rpm -ql oui | grep '/oui$'))
    local path=${paths[0]}
    echo $path
    return 0
}

BIN=$(find_path)

$BIN update
exit 0
