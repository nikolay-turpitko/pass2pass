#!/bin/bash

# This script saves password entry using `pass` comand.

flock .lock pass insert -m "$1"
