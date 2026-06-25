package main

/*
#include <stdio.h>

void helloworld() {
	printf("Hello, World!\n");
}
*/
import "C"

func main() {
	C.helloworld()
}
