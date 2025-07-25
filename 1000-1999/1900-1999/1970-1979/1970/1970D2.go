package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	next   [2]int
	link   int
	length int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b string
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}
	s := a + b
	n := len(s)
	sam := make([]state, 2*n+1)
	last, size := 0, 1
	sam[0].link = -1
	extend := func(ch byte) {
		c := 0
		if ch == 'O' {
			c = 1
		}
		cur := size
		size++
		sam[cur].length = sam[last].length + 1
		p := last
		for p != -1 && sam[p].next[c] == 0 {
			sam[p].next[c] = cur + 1
			p = sam[p].link
		}
		if p == -1 {
			sam[cur].link = 0
		} else {
			q := sam[p].next[c] - 1
			if sam[p].length+1 == sam[q].length {
				sam[cur].link = q
			} else {
				clone := size
				size++
				sam[clone] = sam[q]
				sam[clone].length = sam[p].length + 1
				for p != -1 && sam[p].next[c] == q+1 {
					sam[p].next[c] = clone + 1
					p = sam[p].link
				}
				sam[q].link = clone
				sam[cur].link = clone
			}
		}
		last = cur
	}
	for i := 0; i < n; i++ {
		extend(s[i])
	}
	ans := 0
	for i := 1; i < size; i++ {
		ans += sam[i].length - sam[sam[i].link].length
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
