#!/bin/bash

# This script is used to replace some substring within entry path (effectively
# moving entries to other path during export).

sed -E -e 's%smsaero/taxi%taxi/sms/smsaero%g' | tr -d "\n"
