package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type state struct {
	next [26]int
	link int
	len  int
	pos  []int
}

type SAM struct {
	st   []state
	last int
}

func NewSAM(maxLen int) *SAM {
	st := make([]state, 2*maxLen+5)
	for i := range st {
		for j := 0; j < 26; j++ {
			st[i].next[j] = -1
		}
		st[i].link = -1
	}
	return &SAM{st: st, last: 0}
}

func (sam *SAM) extend(c int, pos int) {
	cur := sam.nextIndex()
	sam.st[cur].len = sam.st[sam.last].len + 1
	sam.st[cur].pos = append(sam.st[cur].pos, pos)
	p := sam.last
	for p >= 0 && sam.st[p].next[c] == -1 {
		sam.st[p].next[c] = cur
		p = sam.st[p].link
	}
	if p == -1 {
		sam.st[cur].link = 0
	} else {
		q := sam.st[p].next[c]
		if sam.st[p].len+1 == sam.st[q].len {
			sam.st[cur].link = q
		} else {
			clone := sam.nextIndex()
			sam.st[clone] = sam.st[q]
			sam.st[clone].len = sam.st[p].len + 1
			sam.st[clone].pos = nil
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

func (sam *SAM) nextIndex() int {
	idx := 0
	for idx < len(sam.st) && sam.st[idx].len != 0 || idx == 0 {
		if sam.st[idx].len == 0 && idx != 0 {
			break
		}
		idx++
	}
	if idx >= len(sam.st) {
		sam.st = append(sam.st, state{})
		for j := 0; j < 26; j++ {
			sam.st[idx].next[j] = -1
		}
		sam.st[idx].link = -1
	} else {
		for j := 0; j < 26; j++ {
			sam.st[idx].next[j] = -1
		}
		sam.st[idx].link = -1
	}
	return idx
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var w string
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &w)
	sam := NewSAM(n)
	pref := make([]int, n+1)
	sam.st[0].len = 0
	sam.st[0].link = -1
	for i := 0; i < n; i++ {
		sam.extend(int(w[i]-'a'), i+1)
		pref[i+1] = sam.last
	}

	size := 0
	for i := range sam.st {
		if sam.st[i].len != 0 || i == 0 {
			size = i + 1
		}
	}
	sam.st = sam.st[:size]

	maxLen := 0
	for i := 0; i < size; i++ {
		if sam.st[i].len > maxLen {
			maxLen = sam.st[i].len
		}
	}
	cnt := make([]int, maxLen+1)
	for i := 0; i < size; i++ {
		cnt[sam.st[i].len]++
	}
	for i := 1; i <= maxLen; i++ {
		cnt[i] += cnt[i-1]
	}
	order := make([]int, size)
	for i := size - 1; i >= 0; i-- {
		l := sam.st[i].len
		cnt[l]--
		order[cnt[l]] = i
	}

	gap := make([]int, size)
	first := make([]int, size)
	const INF = int(1e9)
	for i := 0; i < size; i++ {
		gap[i] = INF
		first[i] = -1
	}

	for _, v := range order {
		if len(sam.st[v].pos) > 1 {
			sort.Ints(sam.st[v].pos)
			best := INF
			idx := 0
			for i := 1; i < len(sam.st[v].pos); i++ {
				d := sam.st[v].pos[i] - sam.st[v].pos[i-1]
				if d < best {
					best = d
					idx = i - 1
				}
			}
			gap[v] = best
			first[v] = sam.st[v].pos[idx]
		}
		if sam.st[v].link >= 0 {
			p := sam.st[v].link
			if len(sam.st[p].pos) < len(sam.st[v].pos) {
				sam.st[p].pos, sam.st[v].pos = sam.st[v].pos, sam.st[p].pos
			}
			sam.st[p].pos = append(sam.st[p].pos, sam.st[v].pos...)
		}
	}

	dp := make([]int, size)
	visited := make([]bool, size)

	var dfs func(int) int
	dfs = func(v int) int {
		if visited[v] {
			return dp[v]
		}
		visited[v] = true
		res := 1
		if gap[v] != INF {
			L := 1
			if sam.st[v].link >= 0 {
				L = sam.st[sam.st[v].link].len + 1
			}
			start := first[v] - L + 1
			end := first[v] + gap[v]
			if start >= 1 && end <= n {
				u := substringState(pref, sam.st, start, end)
				res = 1 + dfs(u)
			}
		}
		dp[v] = res
		return res
	}

	ans := 1
	for i := 1; i < size; i++ {
		if sam.st[i].len > 0 {
			val := dfs(i)
			if val > ans {
				ans = val
			}
		}
	}
	fmt.Println(ans)
}

func substringState(pref []int, st []state, l, r int) int {
	v := pref[r]
	length := r - l + 1
	for st[st[v].link].len >= length {
		v = st[v].link
	}
	return v
}
