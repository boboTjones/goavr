#include "tests.h"

int main() {
  unsigned char buf1[] = "\x01\x02\x03\x04\x05\x06";
  int i = 0;
  int ret = 0;

  memset(buf1, 0, 3);

  for(i = 0; i < sizeof buf1; i++) { 
    ret += buf1[i];    
  }

  return ret & 0xff;
}
