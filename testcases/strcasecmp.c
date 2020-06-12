#include "tests.h"

int main() { 
  char buf1[] = "The Rain In Spain";
  char buf2[] = "the RAIN in SPAIN";
  char buf3[] = "the ra ipse dixit";

  if(strcasecmp(buf1, buf2))  { 
    return 1;
  }

  if(strcasecmp(buf2, buf1))  { 
    return 2;
  }

  if(!strcasecmp(buf2, buf3))  { 
    return 3;
  }

  if(!strcasecmp(buf1, buf3))  { 
    return 4;
  }

  return 5;
}
