char *mystrcpy(char *dest, char *src)
{
    while((*dest++ = *src++)!= '\0')
        ; // <<== Very important!!!
    return dest;
}

char *myitoa(i)
     int i;
{
  static char buf[17];
  char *p = buf + 16;
  if (i >= 0) {
    do {
      *--p = '0' + (i % 10);
      i /= 10;
    } while (i != 0);
    return p;
  }
  else {
    do {
      *--p = '0' - (i % 10);
      i /= 10;
    } while (i != 0);
    *--p = '-';
  }
  return p;
}

int main(void) {
    int foo = 4242;
    char* str;
    char dest[4];
    str = myitoa(foo);
    mystrcpy(dest, str);
    return 0;
}

