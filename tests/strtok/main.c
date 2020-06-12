#include <stdio.h>
#include <string.h>

int
main() { 
  char buf[] = "this:is:a:test";
  char *token = NULL;

  for(token = strtok(buf, ":"); token; token = strtok(NULL, ":"))  {
    puts(token);
  }
  
  return 0;
}
