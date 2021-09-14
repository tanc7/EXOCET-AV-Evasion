package main

/*
#include <stdio.h>
#include <windows.h>

void box()
		{
			MessageBox(0, "Is Go the best?", "C GO GO", 0x00000004L);
	}
*/

import "C"
func main() {
	C.box()
}