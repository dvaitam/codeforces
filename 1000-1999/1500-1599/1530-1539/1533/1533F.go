package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	pref0 := make([]int, n+1)
	pref1 := make([]int, n+1)
	zeroPos := make([]int, 0, n+1)
	onePos := make([]int, 0, n+1)
	for i := 0; i < n; i++ {
		pref0[i+1] = pref0[i]
		pref1[i+1] = pref1[i]
		if s[i] == '0' {
			pref0[i+1]++
			zeroPos = append(zeroPos, i)
		} else {
			pref1[i+1]++
			onePos = append(onePos, i)
		}
	}
	zeroPos = append(zeroPos, n)
	onePos = append(onePos, n)

	ans := make([]int, n)
	for k := 1; k <= n; k++ {
		pos := 0
		cnt := 0
		for pos < n {
			idx0 := pref0[pos] + k
			r0 := n
			if idx0 < len(zeroPos) {
				r0 = zeroPos[idx0]
			}
			idx1 := pref1[pos] + k
			r1 := n
			if idx1 < len(onePos) {
				r1 = onePos[idx1]
			}
			if r1 > r0 {
				r0 = r1
			}
			pos = r0
			cnt++
		}
		ans[k-1] = cnt
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
}
