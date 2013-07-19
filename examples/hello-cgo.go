package main

/*
#include <stdio.h>
void helloworld()
{
  printf("Hello from Cgo!\n");
}
*/
import "C"

func main() {
    C.helloworld()
}
