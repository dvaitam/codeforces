package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]uint64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var ans uint64
	for l := 0; l < n; l++ {
		mn, mx := a[l], a[l]
		for r := l; r < n; r++ {
			if a[r] < mn {
				mn = a[r]
			}
			if a[r] > mx {
				mx = a[r]
			}
			if bits.OnesCount64(mn) == bits.OnesCount64(mx) {
				ans++
			}
		}
	}

	fmt.Fprintln(out, ans)
}
