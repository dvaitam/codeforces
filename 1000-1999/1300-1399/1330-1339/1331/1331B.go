package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}

	for p := 2; p <= a; p++ {
		if a%p == 0 {
			q := a / p
			if isPrime(p) && isPrime(q) {
				fmt.Fprintf(out, "%d%d\n", p, q)
				return
			}
		}
	}
}
