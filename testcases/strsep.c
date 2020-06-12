#include "tests.h"

int main() { 
  char buf1[] = "The Rain In Spain";
  char *cp = buf1, *token = NULL;
  int ret = 0;

  while((token = strsep(&cp, " ")) != NULL) { 
    ret = csum(token);
  }
  
  return ret & 0xff;
}
