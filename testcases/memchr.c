#include "tests.h"

int main() {
  unsigned char buf[] = "\x01\x02\x03\x04\x05\xaa\x06";

  unsigned char *ret = memchr(buf, 0xaa, sizeof(buf));

  if(ret) { 
    return ret[0];
  }

  return 666;
}
