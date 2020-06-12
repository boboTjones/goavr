#include "tests.h"

int main() { 
  int x = 0, ret = 0;

  srand(10);
  rand();
  rand();
  rand();
  x = rand();

  srand(10);
  rand();
  rand();
  rand();
  ret = rand();

  return ret == x;
}
