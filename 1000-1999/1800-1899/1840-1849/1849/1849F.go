package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	type pair struct{ val, idx int }
	b := make([]pair, n)
	for i := range a {
		b[i] = pair{a[i], i}
	}
	sort.Slice(b, func(i, j int) bool { return b[i].val < b[j].val })

	mask := 0
	var ss [][]int
	for i := 29; i >= 0; i-- {
		pre := -1
		cnt := 0
		ok := true
		for j := 0; j < n; j++ {
			curMask := b[j].val & mask
			if curMask != pre {
				pre = curMask
				cnt = 0
			}
			cnt++
			if cnt > 4 {
				ok = false
				break
			}
		}
		if !ok {
			mask |= 1 << i
			continue
		}
		pre = -1
		var cur []int
		for j := 0; j < n; j++ {
			curMask := b[j].val & mask
			if curMask != pre {
				if len(cur) > 0 {
					ss = append(ss, cur)
				}
				pre = curMask
				cur = nil
			}
			cur = append(cur, b[j].idx)
		}
		if len(cur) > 0 {
			ss = append(ss, cur)
		}
		break
	}

	s := make([]byte, n)
	for _, p := range ss {
		l := len(p)
		if l <= 2 {
			for i := 0; i < l; i++ {
				if i&1 == 0 {
					s[p[i]] = '0'
				} else {
					s[p[i]] = '1'
				}
			}
		} else {
			k0, k1 := 0, 1
			if l == 3 {
				if (a[p[0]] ^ a[p[2]]) > (a[p[k0]] ^ a[p[k1]]) {
					k0, k1 = 0, 2
				}
				if (a[p[1]] ^ a[p[2]]) > (a[p[k0]] ^ a[p[k1]]) {
					k0, k1 = 1, 2
				}
				s[p[k0]] = '0'
				s[p[k1]] = '0'
			} else {
				sum := 0
				for _, x := range p {
					sum ^= a[x]
				}
				val := min(a[p[0]]^a[p[1]], sum^(a[p[0]]^a[p[1]]))
				for i := 0; i < 4; i++ {
					for j := 1; j < 4; j++ {
						cur := a[p[i]] ^ a[p[j]]
						mcur := min(cur, sum^cur)
						if mcur > val {
							val = mcur
							k0 = i
							k1 = j
						}
					}
				}
				s[p[k0]] = '0'
				s[p[k1]] = '0'
			}
			for i := 0; i < l; i++ {
				if i != k0 && i != k1 {
					s[p[i]] = '1'
				}
			}
		}
	}
	out.Write(s)
	out.WriteByte('\n')
}
