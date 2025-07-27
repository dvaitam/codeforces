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
		a := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
		}

		// For arrays shorter than 3, no 3-cycle can be applied.
		if n < 3 {
			sorted := true
			for i := 1; i < n; i++ {
				if a[i] < a[i-1] {
					sorted = false
					break
				}
			}
			if sorted {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		// If there is a duplicate value, parity can be adjusted arbitrarily.
		dup := false
		for _, v := range freq {
			if v > 1 {
				dup = true
				break
			}
		}
		if dup {
			fmt.Fprintln(out, "YES")
			continue
		}

		// No duplicates: the array is a permutation of 1..n. Check parity.
		visited := make([]bool, n+1)
		cycles := 0
		for i := 1; i <= n; i++ {
			if !visited[i] {
				cycles++
				for j := i; !visited[j]; j = a[j-1] {
					visited[j] = true
				}
			}
		}
		if (n-cycles)%2 == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
