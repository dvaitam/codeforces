package main

import (
	"bufio"
	"fmt"
	"os"
)

const alpha = 26

type state struct {
	next   [alpha]int
	link   int
	length int
	cnt    int64
}

var (
	st   []state
	size int
	last int
	pos  []int
)

func initSAM(maxLen int) {
	st = make([]state, 2*maxLen+5)
	for i := range st {
		for j := 0; j < alpha; j++ {
			st[i].next[j] = -1
		}
		st[i].link = -1
	}
	size = 1
	last = 0
	st[0].link = -1
	st[0].length = 0
}

func saExtend(c int) {
	cur := size
	size++
	st[cur].length = st[last].length + 1
	for i := 0; i < alpha; i++ {
		st[cur].next[i] = -1
	}
	st[cur].cnt = 0
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
			clone := size
			size++
			st[clone] = st[q]
			st[clone].length = st[p].length + 1
			st[clone].cnt = 0
			for ; p != -1 && st[p].next[c] == q; p = st[p].link {
				st[p].next[c] = clone
			}
			st[q].link = clone
			st[cur].link = clone
		}
	}
	last = cur
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	initSAM(n)
	pos = make([]int, n+1)
	for i := 0; i < n; i++ {
		saExtend(int(s[i] - 'a'))
		pos[i+1] = last
	}

	st = st[:size]
	maxLen := 0
	for i := 0; i < size; i++ {
		if st[i].length > maxLen {
			maxLen = st[i].length
		}
	}

	cntLen := make([]int, maxLen+1)
	for i := 0; i < size; i++ {
		cntLen[st[i].length]++
	}
	for i := 1; i <= maxLen; i++ {
		cntLen[i] += cntLen[i-1]
	}
	order := make([]int, size)
	for i := size - 1; i >= 0; i-- {
		l := st[i].length
		cntLen[l]--
		order[cntLen[l]] = i
	}

	for i := 1; i <= n; i++ {
		if t[i-1] == '0' {
			st[pos[i]].cnt++
		}
	}

	for i := size - 1; i > 0; i-- {
		v := order[i]
		p := st[v].link
		if p >= 0 {
			st[p].cnt += st[v].cnt
		}
	}

	var ans int64
	for i := 1; i < size; i++ {
		if st[i].cnt > 0 {
			val := int64(st[i].length) * st[i].cnt
			if val > ans {
				ans = val
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
