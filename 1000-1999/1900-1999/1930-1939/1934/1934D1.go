package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func solveCase(n, m uint64) ([]uint64, bool) {
	if m >= n || m == 0 {
		return nil, false
	}
	if n&(n-1) == 0 {
		return nil, false
	}
	pow := uint64(1) << (bits.Len64(n) - 1)
	r := n - pow
	if m >= pow {
		// m has the same highest bit as n
		if n^m >= n {
			return nil, false
		}
		return []uint64{n, m}, true
	}
	h := uint64(1) << (bits.Len64(m) - 1)
	if r&h != 0 {
		return []uint64{n, m}, true
	}
	if r > h {
		return []uint64{n, pow + h, m}, true
	}
	return nil, false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m uint64
		fmt.Fscan(in, &n, &m)
		seq, ok := solveCase(n, m)
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, len(seq)-1)
		for i, v := range seq {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
