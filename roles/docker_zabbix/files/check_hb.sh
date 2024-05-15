#!/bin/bash

string=$(cat /synobackup.log | grep -w "\[Local_HyperBackup\] Backup task finished successfully" | tail -n1)
#echo $string
today=$(date +%Y/%m/%d)
#echo $today
date=$(echo $string | awk '{print $2}')
#echo $date
if [[ "$date" == "$today" ]]; then
  echo "1"
else
  echo "2"
fi
