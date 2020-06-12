#include "tests.h"

int main() {
  unsigned char buf1[] = "\x01\x02\x03\x04\x05\xaa\x06";
  unsigned char buf2[] = "\x01\x02\x03\x04\x05\xaa\x06\x77";

  if(memcmp(buf1, buf2, 7)) {
    return 1;
  }

  if(memcmp(buf1, buf2, 3)) {
    return 2;
  }

  if(memcmp(buf1, &buf2[1], 7)) { 
    return 77;
  }  

  return 3;
}
