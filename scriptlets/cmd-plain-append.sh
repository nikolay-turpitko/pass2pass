#!/bin/bash

# This script saves password entry as a plain-text file in the
# output directory. It can be used for debugging.

d="./out"
p="$d/$1"
nd=$(dirname "$p")
mkdir -p "$nd" && cat >> "$p" && echo >> "$p" && echo >> "$p"
