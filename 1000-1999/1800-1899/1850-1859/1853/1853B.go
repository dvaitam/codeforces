package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt. We count valid starting pairs (f1,f2)
// producing the k-th term n in a Fibonacci-like sequence. The k-th
// element equals Fib(k-2)*f1 + Fib(k-1)*f2. Consecutive Fibonacci
// numbers are coprime, so all solutions have the form
// (f1,f2)=(a0+Fib(k-1)*t, b0-Fib(k-2)*t). Using a modular inverse to
// find a0, we simply count non-negative t with f1<=f2.

// extendedEuclid returns x, y, gcd such that ax + by = gcd
func extendedEuclid(a, b int) (int, int, int) {
	if b == 0 {
		return 1, 0, a
	}
	x1, y1, g := extendedEuclid(b, a%b)
	x, y := y1, x1-(a/b)*y1
	return x, y, g
}

func modInverse(a, mod int) int {
	x, _, g := extendedEuclid(a, mod)
	if g != 1 {
		return 0
	}
	if x < 0 {
		x += mod
	}
	return x % mod
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fib := []int{0, 1}
	for fib[len(fib)-1] <= 200000 {
		fib = append(fib, fib[len(fib)-1]+fib[len(fib)-2])
	}

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		if k >= len(fib) || fib[k-1] > n {
			fmt.Fprintln(writer, 0)
			continue
		}
		fk2 := fib[k-2]
		fk1 := fib[k-1]
		fk := fib[k]
		inv := modInverse(fk2, fk1)
		a0 := (n % fk1) * inv % fk1
		b0 := (n - fk2*a0) / fk1
		tMax1 := b0 / fk2
		tMax2 := (b0 - a0) / fk
		if tMax1 < tMax2 {
			tMax2 = tMax1
		}
		if tMax2 < 0 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, tMax2+1)
		}
	}
}
