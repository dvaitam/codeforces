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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		mat := make([][5]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				fmt.Fscan(in, &mat[i][j])
			}
		}

		possible := false
		for d1 := 0; d1 < 5 && !possible; d1++ {
			for d2 := d1 + 1; d2 < 5 && !possible; d2++ {
				aOnly, bOnly := 0, 0
				valid := true
				for i := 0; i < n; i++ {
					a := mat[i][d1]
					b := mat[i][d2]
					if a == 0 && b == 0 {
						valid = false
						break
					}
					if a == 1 && b == 0 {
						aOnly++
					} else if a == 0 && b == 1 {
						bOnly++
					}
				}
				if valid && aOnly <= n/2 && bOnly <= n/2 {
					possible = true
				}
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
