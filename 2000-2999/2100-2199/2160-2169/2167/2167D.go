package main

import (
	"bufio"
	"fmt"
	"os"
)

var primeList = []int64{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
	31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
	73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
	127, 131, 137, 139, 149, 151, 157, 163, 167, 173,
	179, 181, 191, 193, 197, 199, 211, 223, 227, 229,
	233, 239, 241, 251, 257, 263, 269, 271, 277, 281,
	283, 293, 307, 311, 313, 317, 331, 337, 347, 349,
	353, 359, 367, 373, 379, 383, 389, 397, 401, 409,
	419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
	467, 479, 487, 491, 499, 503, 509, 521, 523, 541,
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func isPrime(x int64) bool {
	if x < 2 {
		return false
	}
	if x%2 == 0 {
		return x == 2
	}
	if x%3 == 0 {
		return x == 3
	}
	for i := int64(5); i <= x/i; i += 6 {
		if x%i == 0 || x%(i+2) == 0 {
			return false
		}
	}
	return true
}

func smallestPrimeNotDividing(g int64) int64 {
	if g == 0 {
		return 2
	}
	for _, p := range primeList {
		if g%p != 0 {
			return p
		}
	}
	candidate := primeList[len(primeList)-1] + 2
	if candidate%2 == 0 {
		candidate++
	}
	for {
		if isPrime(candidate) && g%candidate != 0 {
			return candidate
		}
		candidate += 2
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var g int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if i == 0 {
				g = x
			} else {
				g = gcd(g, x)
			}
		}
		fmt.Fprintln(out, smallestPrimeNotDividing(g))
	}
}
