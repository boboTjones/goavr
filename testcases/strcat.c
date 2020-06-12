#include "tests.h"

int main() { 
  char buf1[] = "The Rain In Spain ";
  char buf2[] = "the RAIN in SPAIN";
  char buf3[100] = "";

  strcat(buf3, buf1);
  strcat(buf3, buf2);

  return csum(buf3) & 0xff;
}
