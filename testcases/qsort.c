#include "tests.h"

static int 
compare(const void *a, const void *b) { 
  const unsigned char *ak = a, *bk = b;

  if(*ak < *bk)
    return -1;
  else if(*ak > *bk) 
    return 1;

  return 0;
}

int main() { 
  int ret = 0;

  unsigned char table[] = { 
    0x03,
    0x0c,

    0x02,
    0x0b,

    0x01,
    0x0a,

    0x04,
    0x0d,

    0x06,
    0x0f,

    0x05,
    0x0e,
  };

  qsort(table, 6, 2, compare);

  ret = table[1] + table[11];

  return ret;
}
