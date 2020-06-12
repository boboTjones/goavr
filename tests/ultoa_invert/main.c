#include <stdlib.h>
#include <string.h>
#include <stdio.h>

char * __ultoa_invert (unsigned long val, char * str, int base);

int main(void) {
    char str[16];
    __ultoa_invert(4242, str, 10);
    return 0;
}
