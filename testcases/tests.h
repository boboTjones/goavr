#ifndef TESTS_INCLUDED
#define TESTS_INCLUDED

#include <string.h>
#include <stdlib.h>

unsigned csum(char *key) {
    char *p = key;
    unsigned h = 0, len = strlen(key);
    int i;

    for (i = 0; i < len; i++) {
        h = 33 * h + p[i];
    }

    return h;
}


#endif
