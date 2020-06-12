#include "tests.h"

int main() {
  long i = 10;

  i -= 30;

  i = labs(i);

  return((int)i);
}
