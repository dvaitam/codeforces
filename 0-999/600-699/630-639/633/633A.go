package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b, c int
	if _, err := fmt.Fscan(os.Stdin, &a, &b, &c); err != nil {
		return
	}
	for x := 0; a*x <= c; x++ {
		if (c-a*x)%b == 0 {
			fmt.Println("Yes")
			return
		}
	}
	fmt.Println("No")
}
