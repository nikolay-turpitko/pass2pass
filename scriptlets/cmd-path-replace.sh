#!/bin/bash

# This script is used to replace some substring within entry path (effectively
# moving entries to other path during export).

sed -E \
    -e 's%^private/bank/tinkoff.ru%BANK/tinkoff/private%g' \
    -e 's%^business/bank/tinkoff.ru%BANK/tinkoff/business%g' \
    -e 's%^private/bank%depricated/bank%g' \
    -e 's%^business/bank%depricated/bank%g' \
    -e 's%^business/social%social%g' \
    -e 's%^private/pay-online%BANK/PAY-ONLINE%g' \
    -e 's%^private/telephony%telephony%g' \
    -e 's%^private/net%NET%g' \
    -e 's%^private/other/(gosuslugi|nalog|госуслуги|pochta|sistemagorod)%STATE/\1%g' \
    -e 's%^private/other/(.*taximaxim|s7.*)%transport/\1%g' \
    -e 's%^private/other/(.*asus|.*\.hp|zyxel)%hardware/\1%g' \
    -e 's%^private/other/hightail%depricated/hightail%g' \
    -e 's%^private/mail-and-im/(icq|sibnet|kate.*)%depricated/email-and-im/\1%g' \
    -e 's%^private/mail-and-im/mail.ru%email-and-im/mailru/personal%g' \
    -e 's%^private/mail-and-im/.*disposal.*/%email-and-im/google/disposal/%g' \
    -e 's%^private/mail-and-im/google/.*/%email-and-im/google/personal/%g' \
    -e 's%^private/mail-and-im/zoho/.*/%email-and-im/zoho/personal/%g' \
    -e 's%^private/mail-and-im/bitlbee%email-and-im/bitlbee%g' \
    -e 's%^private/other/(.*)/%JOB/\1/%g' \
    -e 's%^private%%g' \
    -e 's%^business/job-hunting%JOB/JOB-HUNTING%g' \
    -e 's%^business/proffessional%education%g' \
    -e 's%^business/mail-and-im/(aim|.*\.aol|yahoo!)%depricated/email-and-im/\1%g' \
    -e 's%^business/mail-and-im/godaddy%JOB/godaddy%g' \
    -e 's%^business/mail-and-im/google/.*/%email-and-im/google/work/%g' \
    -e 's%^business/mail-and-im/zoho/.*/%email-and-im/zoho/work/%g' \
    -e 's%^business/mail-and-im/zoho.eu%email-and-im/zoho/zoho.eu%g' \
    -e 's%^business/mail-and-im/freenode%email-and-im/freenode%g' \
    -e 's%^business/work/galileo/trello%JOB/trello%g' \
    -e 's%^business/work/ponominalu%JOB/letsrock/ponominalu%g' \
    -e 's%^business/work/seatgeek%JOB/letsrock/seatgeek%g' \
    -e 's%^business/work%JOB%g' \
    -e 's%^business/other/(smsaero|twilio)%JOB/\1%g' \
    -e 's%^business/other/(paypal-tech-support)%BANK/PAY-ONLINE/\1%g' \
    -e 's%^business/other/(elba.kontur)%JOB/BIZ/\1%g' \
    -e 's%JOB/letsrock%JOB/LETSROCK%g' \
    -e 's%JOB/galileo%JOB/GALILEO%g' \
    -e 's%http.?-%%g' \
    -e 's%(\.ru)|(\.com)|(\.net)|(\.org)%%g' | tr -d "\n"
