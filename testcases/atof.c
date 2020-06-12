#include "tests.h"

int main() {
  double ret;
  char *str = "1024.2";

  ret = atof(str);

  ret *= 10;

  return (int)ret;

}
