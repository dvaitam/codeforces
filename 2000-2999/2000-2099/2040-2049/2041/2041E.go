package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}

	if a == b {
		fmt.Println(1)
		fmt.Println(a)
		return
	}

	x := 3*a - 2*b
	fmt.Println(3)
	fmt.Printf("%d %d %d\n", b, b, x)
}
