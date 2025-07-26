package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1881B - Three Threadlets.
// We try all possible final counts of threadlets (3 to 6).
// For each count k dividing the total length, every initial
// thread length must also be divisible by the target piece length.
// Since k <= 6, the required number of cuts k-3 never exceeds 3.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, c int64
		fmt.Fscan(in, &a, &b, &c)
		sum := a + b + c
		ok := false
		for k := int64(3); k <= 6; k++ {
			if sum%k != 0 {
				continue
			}
			l := sum / k
			if a%l == 0 && b%l == 0 && c%l == 0 {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
