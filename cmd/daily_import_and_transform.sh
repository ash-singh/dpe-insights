#! /bin/sh

printf "\nRunning GitHub Daily cron\n"
./build/dpe-cli github cron --days=1

printf "\nImport incidents from PagerDuty\n"
./build/dpe-cli pagerduty cron --days=14
