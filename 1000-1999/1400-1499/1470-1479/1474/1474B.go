package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(x int) int {
	for {
		if isPrime(x) {
			return x
		}
		x++
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var d int
		fmt.Fscan(reader, &d)
		p1 := nextPrime(1 + d)
		p2 := nextPrime(p1 + d)
		fmt.Fprintln(writer, p1*p2)
	}
}
