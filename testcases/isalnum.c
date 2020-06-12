#include "tests.h"
#include <ctype.h>

int main() { 
  char buf[] = "a(*$98nFN89$*(nnc48$*CNm4c*C$M*";

  int ret = 0;
  int i = 0;

  for(i = 0; i < sizeof buf; i++) { 
    if(isalnum(buf[i])) 
      ret += 1;
  }

  return ret;
}
