package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type maskInfo struct {
	mask int
	pop  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var d, n int
	fmt.Fscan(in, &d, &n)

	info := make([]maskInfo, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		mask := 0
		for j := 0; j < d; j++ {
			if s[j] == '1' {
				mask |= 1 << j
			}
		}
		info[i] = maskInfo{mask: mask, pop: bits.OnesCount(uint(mask))}
	}

	sort.Slice(info, func(i, j int) bool {
		if info[i].pop == info[j].pop {
			return info[i].mask < info[j].mask
		}
		return info[i].pop < info[j].pop
	})

	var seq []string
	cur := 0
	for _, node := range info {
		mask := node.mask
		if cur&mask != cur {
			if cur != 0 {
				seq = append(seq, "R")
			}
			cur = 0
		}
		diff := mask &^ cur
		for diff > 0 {
			lsb := diff & -diff
			idx := bits.TrailingZeros(uint(lsb))
			seq = append(seq, fmt.Sprintf("%d", idx))
			cur |= lsb
			diff &= diff - 1
		}
	}

	fmt.Fprintln(out, len(seq))
	for i, token := range seq {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, token)
	}
	fmt.Fprintln(out)
}
