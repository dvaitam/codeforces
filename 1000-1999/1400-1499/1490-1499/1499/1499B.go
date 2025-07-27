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
		var s string
		fmt.Fscan(in, &s)
		n := len(s)

		dp := make([][][][]bool, n+1)
		for i := range dp {
			dp[i] = make([][][]bool, 2)
			for p := range dp[i] {
				dp[i][p] = make([][]bool, 2)
				for r := range dp[i][p] {
					dp[i][p][r] = make([]bool, 2)
				}
			}
		}
		dp[0][0][0][0] = true
		for i := 0; i < n; i++ {
			ch := s[i]
			for phase := 0; phase < 2; phase++ {
				for prev := 0; prev < 2; prev++ {
					for rem := 0; rem < 2; rem++ {
						if !dp[i][phase][prev][rem] {
							continue
						}
						if phase == 0 {
							if ch == '0' {
								dp[i+1][0][0][rem] = true
							} else {
								dp[i+1][1][0][rem] = true
							}
						} else {
							if ch == '1' {
								dp[i+1][1][0][rem] = true
							}
						}
						if prev == 0 {
							dp[i+1][phase][1][rem|1] = true
						}
					}
				}
			}
		}
		ok := false
		for phase := 0; phase < 2; phase++ {
			for prev := 0; prev < 2; prev++ {
				if dp[n][phase][prev][1] {
					ok = true
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
