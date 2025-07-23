package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	next   [26]int
	link   int
	length int
	occ    int64
}

type SAM struct {
	st   []state
	last int
	size int
}

func NewSAM(maxLen int) *SAM {
	st := make([]state, 2*maxLen+5)
	for i := range st {
		for j := 0; j < 26; j++ {
			st[i].next[j] = -1
		}
		st[i].link = -1
	}
	st[0].link = -1
	return &SAM{st: st, last: 0, size: 1}
}

func (sam *SAM) Extend(c int) {
	cur := sam.size
	sam.size++
	sam.st[cur].length = sam.st[sam.last].length + 1
	sam.st[cur].occ = 1
	for j := 0; j < 26; j++ {
		sam.st[cur].next[j] = -1
	}
	p := sam.last
	for p >= 0 && sam.st[p].next[c] == -1 {
		sam.st[p].next[c] = cur
		p = sam.st[p].link
	}
	if p == -1 {
		sam.st[cur].link = 0
	} else {
		q := sam.st[p].next[c]
		if sam.st[p].length+1 == sam.st[q].length {
			sam.st[cur].link = q
		} else {
			clone := sam.size
			sam.size++
			sam.st[clone] = sam.st[q]
			sam.st[clone].length = sam.st[p].length + 1
			sam.st[clone].occ = 0
			for p >= 0 && sam.st[p].next[c] == q {
				sam.st[p].next[c] = clone
				p = sam.st[p].link
			}
			sam.st[q].link = clone
			sam.st[cur].link = clone
		}
	}
	sam.last = cur
}

func (sam *SAM) Prepare() []int {
	tot := sam.size
	maxLen := 0
	for i := 0; i < tot; i++ {
		if sam.st[i].length > maxLen {
			maxLen = sam.st[i].length
		}
	}
	cnt := make([]int, maxLen+1)
	for i := 0; i < tot; i++ {
		cnt[sam.st[i].length]++
	}
	for i := 1; i <= maxLen; i++ {
		cnt[i] += cnt[i-1]
	}
	order := make([]int, tot)
	for i := tot - 1; i >= 0; i-- {
		l := sam.st[i].length
		cnt[l]--
		order[cnt[l]] = i
	}
	for i := tot - 1; i > 0; i-- {
		v := order[i]
		p := sam.st[v].link
		if p >= 0 {
			sam.st[p].occ += sam.st[v].occ
		}
	}
	return order
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		sam := NewSAM(len(s))
		for i := 0; i < len(s); i++ {
			sam.Extend(int(s[i] - 'a'))
		}
		sam.Prepare()
		var result int64
		for i := 1; i < sam.size; i++ {
			v := &sam.st[i]
			linkLen := 0
			if v.link >= 0 {
				linkLen = sam.st[v.link].length
			}
			lengthDiff := v.length - linkLen
			result += v.occ * v.occ * int64(lengthDiff)
		}
		fmt.Fprintln(out, result)
	}
}
