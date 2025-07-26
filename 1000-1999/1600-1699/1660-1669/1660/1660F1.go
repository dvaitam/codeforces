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
		var s string
		fmt.Fscan(in, &s)
		ans := 0
		for i := 0; i < n; i++ {
			diff := 0  // minus minus plus difference
			pairs := 0 // maximum pairs of adjacent '-'
			run := 0   // current run length of '-'
			for j := i; j < n; j++ {
				if s[j] == '-' {
					diff++
					run++
					if run%2 == 0 {
						pairs++
					}
				} else {
					diff--
					run = 0
				}
				if diff >= 0 && diff%3 == 0 && pairs >= diff/3 {
					ans++
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
