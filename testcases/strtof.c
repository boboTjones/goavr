#include "tests.h"

int main() {
  float ret;
  char *str = "-1024.2yabbadabba";
  char *ep = NULL;

  ret = strtof(str, &ep);

  if(ep != NULL && *ep == 'y' && (int)ret == -(1024))
    return 123; 

  return 0;
}
