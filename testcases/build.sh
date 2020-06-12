#!/bin/bash

for i in `ls *.c | awk -F"." '{print$1}'`
do
  echo "Compiling $i.c"
  avr-gcc -Wall -Os -DF_CPU=8000000 -mmcu=atmega8 -c $i.c -o $i.o
  avr-gcc -Wall -Os -DF_CPU=8000000 -mmcu=atmega8 -o $i.elf $i.o
done
