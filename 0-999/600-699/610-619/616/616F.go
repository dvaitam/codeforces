package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const ALPH = 26

type node struct {
	next   [ALPH]int
	link   int
	length int
	val    int64
}

var sam []node
var last, sz int

func samInit(maxLen int) {
	sam = make([]node, 2*maxLen+5)
	for i := range sam {
		for j := 0; j < ALPH; j++ {
			sam[i].next[j] = -1
		}
	}
	last = 0
	sz = 1
	sam[0].link = -1
	sam[0].length = 0
}

func samExtend(ch int, addVal int64) {
	cur := sz
	sz++
	sam[cur].length = sam[last].length + 1
	for j := 0; j < ALPH; j++ {
		sam[cur].next[j] = -1
	}
	sam[cur].val = addVal
	p := last
	for p != -1 && sam[p].next[ch] == -1 {
		sam[p].next[ch] = cur
		p = sam[p].link
	}
	if p == -1 {
		sam[cur].link = 0
	} else {
		q := sam[p].next[ch]
		if sam[p].length+1 == sam[q].length {
			sam[cur].link = q
		} else {
			clone := sz
			sz++
			sam[clone] = sam[q]
			sam[clone].length = sam[p].length + 1
			sam[clone].val = 0
			for p != -1 && sam[p].next[ch] == q {
				sam[p].next[ch] = clone
				p = sam[p].link
			}
			sam[q].link = clone
			sam[cur].link = clone
		}
	}
	last = cur
}

func propagate() {
	maxL := 0
	for i := 0; i < sz; i++ {
		if sam[i].length > maxL {
			maxL = sam[i].length
		}
	}
	cnt := make([]int, maxL+1)
	for i := 0; i < sz; i++ {
		cnt[sam[i].length]++
	}
	for i := 1; i <= maxL; i++ {
		cnt[i] += cnt[i-1]
	}
	order := make([]int, sz)
	for i := sz - 1; i >= 0; i-- {
		l := sam[i].length
		cnt[l]--
		order[cnt[l]] = i
	}
	for i := sz - 1; i > 0; i-- {
		v := order[i]
		p := sam[v].link
		if p >= 0 {
			sam[p].val += sam[v].val
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
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
	samInit(total)
	for i := 0; i < n; i++ {
		last = 0
		for j := 0; j < len(strs[i]); j++ {
			samExtend(int(strs[i][j]-'a'), costs[i])
		}
	}
	propagate()
	ans := int64(math.MinInt64)
	for i := 1; i < sz; i++ {
		v := sam[i].val * int64(sam[i].length)
		if v > ans {
			ans = v
		}
	}
	fmt.Fprintln(out, ans)
}
