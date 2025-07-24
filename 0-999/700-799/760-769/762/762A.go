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

	var n int64
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	small := make([]int64, 0)
	big := make([]int64, 0)
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			small = append(small, i)
			if i != n/i {
				big = append(big, n/i)
			}
		}
	}

	total := int64(len(small) + len(big))
	if k > total {
		fmt.Fprintln(out, -1)
		return
	}

	if k <= int64(len(small)) {
		fmt.Fprintln(out, small[k-1])
	} else {
		idx := int64(len(big)) - (k - int64(len(small)))
		fmt.Fprintln(out, big[idx])
	}
}
