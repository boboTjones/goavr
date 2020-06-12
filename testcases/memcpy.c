#include "tests.h"

int main() {
  unsigned char buf1[] = "\x01\x02\x03\x04\x05\xaa\x06";
  unsigned char buf2[] = "\x01\x02\x03\x04\x05\xaa\x06\x77";

  memcmp(&buf2[1], buf1, 7);

  return buf2[1];
}
