#!/bin/bash

# This script saves password entry using `pass` comand.
# Mac OS X: https://github.com/discoteq/flock

flock .lock gopass insert -m "$1"
