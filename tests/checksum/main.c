#include <stdio.h>

int
checksum(char *s) {
  int cs = 0;

  while(*s++)
    cs += *s;

  return cs;
}

int
main() {
  return checksum("yellow sumbarine");
}
