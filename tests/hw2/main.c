#include <avr/io.h>
#include <stdio.h>
#include <string.h>

int main()
{
    char string[] = "Hello, World!";
    int d = strlen(string);
    printf("%s %d\n", string, d);
    return 0;
}
