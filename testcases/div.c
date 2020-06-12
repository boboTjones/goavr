#include "tests.h"

int main() { 
  int num = 53;
  int den = 3;

  div_t res = div(num, den);

  return res.quot;
}
