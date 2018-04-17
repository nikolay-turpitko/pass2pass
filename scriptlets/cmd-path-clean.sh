#!/bin/bash

# This script is used to cleanup fields (Name & Grouping) before use them
# to build entry path (replacing $name & $group in template file path).

tr "[:upper:]" "[:lower:]" | \
    tr -s "\\" "/" | \
    sed -E \
        -e 's%[[:space:]]+-[[:space:]]+%/%g' \
        -e 's%:/+%-%g' \
        -e 's%[[:space:]]+%-%g' \
        -e 's%&amp;%and%g' | \
    tr -d "\n"
