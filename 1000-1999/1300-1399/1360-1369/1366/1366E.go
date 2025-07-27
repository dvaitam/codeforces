package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}

	idx := n - 1
	ans := int64(1)
	for i := m - 1; i >= 0; i-- {
		last := -1
		for idx >= 0 && a[idx] >= b[i] {
			if a[idx] == b[i] && last == -1 {
				last = idx
			}
			idx--
		}
		if last == -1 {
			fmt.Fprintln(writer, 0)
			return
		}
		if i == 0 {
			if idx != -1 {
				fmt.Fprintln(writer, 0)
				return
			}
		} else {
			ans = ans * int64(last-idx) % mod
		}
	}
	fmt.Fprintln(writer, ans%mod)
}
