package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		solve(n, k, s, writer)
	}
}

func solve(n, k int, s string, writer *bufio.Writer) {
	m := k
	if m > 20 {
		m = 20
	}
	size := 1 << m
	seen := make([]bool, size)

	// prefix of zeros to check if prefix part is all ones
	zeros := make([]int, n+1)
	for i := 0; i < n; i++ {
		zeros[i+1] = zeros[i]
		if s[i] == '0' {
			zeros[i+1]++
		}
	}

	mask := 0
	for i := 0; i < n; i++ {
		mask = ((mask << 1) & (size - 1))
		if s[i] == '1' {
			mask |= 1
		}
		if i >= k-1 {
			start := i - k + 1
			if k > m {
				if zeros[start+k-m]-zeros[start] == 0 {
					seen[mask] = true
				}
			} else {
				seen[mask] = true
			}
		}
	}

	ans := -1
	for y := 0; y < size; y++ {
		var idx int
		if k > m {
			idx = (^y) & (size - 1)
		} else {
			idx = (size - 1) ^ y
		}
		if !seen[idx] {
			ans = y
			break
		}
	}

	if ans == -1 {
		fmt.Fprintln(writer, "NO")
		return
	}

	fmt.Fprintln(writer, "YES")
	if k > m {
		for i := 0; i < k-m; i++ {
			fmt.Fprint(writer, "0")
		}
	}
	fmt.Fprintf(writer, "%0*b\n", m, ans)
}
