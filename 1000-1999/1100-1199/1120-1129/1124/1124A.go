package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	fmt.Fscan(in, &a, &b)
	if gcd(a, b) > 1 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
