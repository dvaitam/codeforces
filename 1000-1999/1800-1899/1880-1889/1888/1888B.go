package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int64
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}
	fmt.Println(gcd(a, b))
}

