#!/bin/bash

# This script is used to filter out some entries by their path.

grep -v \
    -e 'business/work.*/taxi' \
    -e 'beeline' \
    -e 'appengine' \
    -e 'play.google' \
    -e 'logmein-hamachi' \
    -e 'ronte' || true
