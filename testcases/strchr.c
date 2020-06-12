#include "tests.h"

int main() {
  char buf[] = "abcdefXABCDEF";
  char *ret = NULL;

  if((ret = strchr(buf, 'X'))) { 
    return *ret;
  }

  return 1;
}
