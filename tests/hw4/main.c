#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main()
{
    FILE *nstdOut;
    fprintf(nstdOut, "Hello, World!\n");
    fclose(nstdOut);
    return (0);
}
