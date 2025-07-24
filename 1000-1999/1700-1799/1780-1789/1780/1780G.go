package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// state represents a suffix automaton state
type state struct {
	next   [26]int
	link   int
	length int
	occ    int
}

var st []state
var last int
var sz int

func saInit(n int) {
	st = make([]state, 2*n+5)
	for i := range st {
		for j := 0; j < 26; j++ {
			st[i].next[j] = -1
		}
	}
	st[0].link = -1
	st[0].length = 0
	last = 0
	sz = 1
}

func saExtend(c int) {
	cur := sz
	sz++
	st[cur].length = st[last].length + 1
	st[cur].occ = 1
	for j := 0; j < 26; j++ {
		st[cur].next[j] = -1
	}
	p := last
	for p != -1 && st[p].next[c] == -1 {
		st[p].next[c] = cur
		p = st[p].link
	}
	if p == -1 {
		st[cur].link = 0
	} else {
		q := st[p].next[c]
		if st[p].length+1 == st[q].length {
			st[cur].link = q
		} else {
			clone := sz
			sz++
			st[clone] = st[q]
			st[clone].length = st[p].length + 1
			st[clone].occ = 0
			for p != -1 && st[p].next[c] == q {
				st[p].next[c] = clone
				p = st[p].link
			}
			st[q].link = clone
			st[cur].link = clone
		}
	}
	last = cur
}

func intSqrt(x int) int {
	r := int(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	saInit(n)
	for i := 0; i < n; i++ {
		saExtend(int(s[i] - 'a'))
	}

	// propagate occurrence counts
	maxLen := 0
	for i := 0; i < sz; i++ {
		if st[i].length > maxLen {
			maxLen = st[i].length
		}
	}
	cnt := make([]int, maxLen+1)
	for i := 0; i < sz; i++ {
		cnt[st[i].length]++
	}
	for i := 1; i <= maxLen; i++ {
		cnt[i] += cnt[i-1]
	}
	order := make([]int, sz)
	for i := sz - 1; i >= 0; i-- {
		l := st[i].length
		cnt[l]--
		order[cnt[l]] = i
	}
	for i := sz - 1; i > 0; i-- {
		v := order[i]
		p := st[v].link
		if p >= 0 {
			st[p].occ += st[v].occ
		}
	}

	ans := int64(0)
	for v := 1; v < sz; v++ {
		occ := st[v].occ
		if occ == 0 {
			continue
		}
		l := 1
		if st[v].link >= 0 {
			l = st[st[v].link].length + 1
		}
		r := st[v].length
		limit := intSqrt(occ)
		for d := 1; d <= limit; d++ {
			if occ%d == 0 {
				if d >= l && d <= r {
					ans += int64(occ)
				}
				other := occ / d
				if other != d && other >= l && other <= r {
					ans += int64(occ)
				}
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
