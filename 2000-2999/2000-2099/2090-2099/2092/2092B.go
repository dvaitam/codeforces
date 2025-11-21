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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)

		ones1 := 0
		ones2 := 0
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				if a[i] == '1' {
					ones1++
				}
				if b[i] == '1' {
					ones2++
				}
			} else {
				if b[i] == '1' {
					ones1++
				}
				if a[i] == '1' {
					ones2++
				}
			}
		}
		bCount1 := n / 2
		bCount2 := (n + 1) / 2
		if ones1 <= bCount1 && ones2 <= bCount2 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
