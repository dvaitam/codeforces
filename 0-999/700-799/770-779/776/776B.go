package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	if n <= 2 {
		fmt.Fprintln(out, 1)
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteByte('1')
		}
		if n > 0 {
			out.WriteByte('\n')
		}
		return
	}

	// Use two colors: 1 for primes, 2 for composites
	limit := n + 1
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}

	fmt.Fprintln(out, 2)
	for i := 2; i <= limit; i++ {
		if i > 2 {
			out.WriteByte(' ')
		}
		if isPrime[i] {
			out.WriteByte('1')
		} else {
			out.WriteByte('2')
		}
	}
	out.WriteByte('\n')
}
