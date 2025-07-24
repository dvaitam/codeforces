package main

import (
	"bufio"
	"fmt"
	"os"
)

// The gcd of A! and B! for A,B up to 1e9 with min(A,B)<=12 is simply
// the factorial of the smaller of A and B because (min(A,B))! divides
// both A! and B! and any larger factor would exceed min(A,B)!. Therefore
// we only need to compute factorial of the smaller number (<=12), which
// easily fits into 64-bit integer.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	n := a
	if b < n {
		n = b
	}
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	fmt.Println(res)
}
