#include "tests.h"

int main() {
  long ret = 0;
  char *str = "12345";

  ret = atol(str);
  ret &= 0xff;

  return ret;
}
