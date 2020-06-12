#include <stdio.h>
#include <stdlib.h>

int main()
{
    char s[16];
    sprintf(s, "Hello, World!\n");
    int x = atoi(&s[4]);
    return x;
}
