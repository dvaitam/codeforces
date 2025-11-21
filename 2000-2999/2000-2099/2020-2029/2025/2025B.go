package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	n := make([]int, t)
	k := make([]int, t)
	maxK := 0

	for i := 0; i < t; i++ {
		fmt.Fscan(in, &n[i])
	}
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &k[i])
		if k[i] > maxK {
			maxK = k[i]
		}
	}

	pows := make([]int, maxK+1)
	pows[0] = 1
	for i := 1; i <= maxK; i++ {
		pows[i] = (pows[i-1] * 2) % mod
	}

	for i := 0; i < t; i++ {
		if i > 0 {
			out.WriteByte('\n')
		}
		fmt.Fprint(out, pows[k[i]])
	}
}
