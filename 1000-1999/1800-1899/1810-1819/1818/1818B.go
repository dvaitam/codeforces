package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt from contest 1818.
// If n>1 is odd, the sum of the whole permutation n(n+1)/2 is divisible by n,
// so no arrangement works. For even n, the permutation 2 1 4 3 ... ensures
// every subarray of length >=2 has an odd sum when length is 2 and changing
// parity prevents divisibility for longer lengths. For n=1, output 1.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n%2 == 1 && n > 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		if n == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		for i := 1; i <= n; i += 2 {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, i+1, " ", i)
		}
		fmt.Fprintln(out)
	}
}
