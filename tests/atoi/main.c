#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    char *s = "4242";
    int foo = atoi(s);    
    
    char str[16];
    sprintf(str, "Hello, %d!\n", foo);
    return 0;
}
