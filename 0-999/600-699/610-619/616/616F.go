package main

import (
	"bufio"
	"fmt"
	"os"
)

// Generalized suffix automaton approach for 616F
// For each string s, f(s) = |s| * sum(c_i * p_{s,i})
// where p_{s,i} = number of occurrences of s in t_i.
// We want max f(s) over all strings s.

const ALPH = 26

type saNode struct {
	ch   [ALPH]int32
	link int32
	len  int32
	cnt  int64 // will hold sum of c_i * (number of endpos in string i)
}

var (
	sa   []saNode
	sz   int32
	last int32
)

func saInit(cap int) {
	sa = make([]saNode, 0, 2*cap+10)
	sa = append(sa, saNode{})
	for i := 0; i < ALPH; i++ {
		sa[0].ch[i] = -1
	}
	sa[0].link = -1
	sz = 1
	last = 0
}

func saExtend(c int, val int64) {
	// Check if transition from last already exists (generalized SAM)
	if sa[last].ch[c] != -1 {
		q := sa[last].ch[c]
		if sa[last].len+1 == sa[q].len {
			sa[q].cnt += val
			last = q
			return
		}
		// Clone q
		clone := sz
		sz++
		sa = append(sa, sa[q])
		sa[clone].len = sa[last].len + 1
		sa[clone].cnt = val
		// Redirect
		p := last
		for p != -1 && sa[p].ch[c] == q {
			sa[p].ch[c] = clone
			p = sa[p].link
		}
		sa[q].link = clone
		last = clone
		return
	}

	cur := sz
	sz++
	nd := saNode{}
	for i := 0; i < ALPH; i++ {
		nd.ch[i] = -1
	}
	nd.len = sa[last].len + 1
	nd.cnt = val
	sa = append(sa, nd)

	p := last
	for p != -1 && sa[p].ch[c] == -1 {
		sa[p].ch[c] = cur
		p = sa[p].link
	}
	if p == -1 {
		sa[cur].link = 0
	} else {
		q := sa[p].ch[c]
		if sa[p].len+1 == sa[q].len {
			sa[cur].link = q
		} else {
			clone := sz
			sz++
			cnd := sa[q]
			cnd.len = sa[p].len + 1
			cnd.cnt = 0
			sa = append(sa, cnd)
			for p != -1 && sa[p].ch[c] == q {
				sa[p].ch[c] = clone
				p = sa[p].link
			}
			sa[q].link = clone
			sa[cur].link = clone
		}
	}
	last = cur
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	strs := make([]string, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &strs[i])
		total += len(strs[i])
	}
	costs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &costs[i])
	}

	saInit(total)
	for i := 0; i < n; i++ {
		last = 0
		for j := 0; j < len(strs[i]); j++ {
			saExtend(int(strs[i][j]-'a'), costs[i])
		}
	}

	// Topological sort by length, then propagate cnt to parent via link
	maxLen := int32(0)
	for i := int32(0); i < sz; i++ {
		if sa[i].len > maxLen {
			maxLen = sa[i].len
		}
	}
	bucket := make([]int32, maxLen+1)
	for i := int32(0); i < sz; i++ {
		bucket[sa[i].len]++
	}
	for i := int32(1); i <= maxLen; i++ {
		bucket[i] += bucket[i-1]
	}
	order := make([]int32, sz)
	for i := sz - 1; i >= 0; i-- {
		l := sa[i].len
		bucket[l]--
		order[bucket[l]] = i
	}
	for i := sz - 1; i > 0; i-- {
		v := order[i]
		p := sa[v].link
		if p >= 0 {
			sa[p].cnt += sa[v].cnt
		}
	}

	ans := int64(0)
	for i := int32(1); i < sz; i++ {
		v := sa[i].cnt * int64(sa[i].len)
		if v > ans {
			ans = v
		}
	}
	fmt.Fprintln(out, ans)
}
