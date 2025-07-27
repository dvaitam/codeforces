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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
	}

	var ans int64
	for i := 0; i < n; i += 2 {
		bal := int64(0)
		mn := int64(0)
		for j := i + 1; j < n; j++ {
			if j%2 == 1 { // closing segment
				L := int64(1)
				if -mn > L-1 {
					L = -mn
				}
				if 1-bal > L {
					L = 1 - bal
				}
				R := c[i]
				if c[j]-bal < R {
					R = c[j] - bal
				}
				if L <= R {
					ans += R - L + 1
				}
				bal -= c[j]
				if bal < mn {
					mn = bal
				}
			} else {
				bal += c[j]
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
