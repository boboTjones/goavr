#include <avr/io.h>
#include<stdio.h>

int main(void)
{
    char string[] = "Hello, World!";
    int d = strlen(string);
    printf("%s %d\n", string, d);
    return(1);
}
