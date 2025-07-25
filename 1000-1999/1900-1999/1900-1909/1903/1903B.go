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
		M := make([][]int, n)
		for i := 0; i < n; i++ {
			M[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &M[i][j])
			}
		}

		a := make([]int, n)
		const allOnes = (1 << 30) - 1
		for i := 0; i < n; i++ {
			x := allOnes
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				x &= M[i][j]
			}
			// choose 0 when there is only one element
			if n == 1 {
				x = 0
			}
			a[i] = x
		}

		valid := true
		for i := 0; i < n && valid; i++ {
			for j := i + 1; j < n; j++ {
				if (a[i] | a[j]) != M[i][j] {
					valid = false
					break
				}
			}
		}
		if !valid {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			for i := 0; i < n; i++ {
				if i > 0 {
					out.WriteByte(' ')
				}
				fmt.Fprint(out, a[i])
			}
			fmt.Fprintln(out)
		}
	}
}
