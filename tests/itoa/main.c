#include <stdlib.h>
#include <string.h>

int main(void) {
    int foo = 4242;
    char xxx;    
    char dest[16];
    xxx = *itoa(foo, dest, 10);
    char str[16];
    strcpy(str, dest);
    return 0;
}

