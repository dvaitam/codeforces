package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solve(x, d int64) bool {
	k := 0
	tmp := x
	for tmp%d == 0 {
		tmp /= d
		k++
	}
	r := tmp
	if k <= 1 {
		return false
	}
	if r != 1 && !isPrime(r) {
		return true
	}
	if k <= 2 {
		return false
	}
	if isPrime(d) {
		return false
	}
	p := int64(math.Sqrt(float64(d)))
	if p*p == d && isPrime(p) && r == p && k == 3 {
		return false
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var x, d int64
		fmt.Fscan(reader, &x, &d)
		if solve(x, d) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
