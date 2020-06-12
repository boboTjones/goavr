#!/bin/bash

GOPATH=/Users/erin/codebase/fouravr

for i in `ls *.elf | awk -F"." '{print$1}'`
do 
  RES=`$GOPATH/bin/fouravr -f $i.elf`
  REAL=`cat $i.ret`
  if [ $RES -eq $REAL ]
   then
    echo "$i passed"
   else
    echo "$i failed. Got $RES, should have been $REAL."
  fi
done
