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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if n == 1 {
		fmt.Fprintln(out, a[0])
		return
	}

	even := make([]int64, 2*n+1)
	odd := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		v := a[i%n]
		if i%2 == 0 {
			even[i+1] = even[i] + v
			odd[i+1] = odd[i]
		} else {
			odd[i+1] = odd[i] + v
			even[i+1] = even[i]
		}
	}

	k := n / 2
	var best int64
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			s := even[i+2*k+1] - even[i]
			if s > best {
				best = s
			}
		} else {
			s := odd[i+2*k+1] - odd[i]
			if s > best {
				best = s
			}
		}
	}
	fmt.Fprintln(out, best)
}
