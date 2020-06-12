#include <stdlib.h>
#include <string.h>
#include <stdio.h>

int main(void) {
    char str[24];
    char ins[6] = "world";
    sprintf(str, "Hello, %s!", ins);
    return 0;
}
