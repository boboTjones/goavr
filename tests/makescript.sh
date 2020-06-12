#!/bin/bash 

BASEPATH=`pwd`/tests

for i in `echo $BASEPATH/*`
do 
  if [ -d $i ]
    then cd $i
    make
    cd ../
   fi
done
