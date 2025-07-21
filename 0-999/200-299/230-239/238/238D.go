package main

import (
	"bufio"
	"fmt"
	"os"
)

type State struct {
	exitSide int // 0=left,1=right
	exitDir  int // 0=left,1=right
	cnt      [10]int
}

// Node holds behavior for segment: entry 0 = enter left with DP=right; entry 1 = enter right with DP=left
type Node struct {
	f [2]State
}

func makeNode(c byte) Node {
	var n Node
	if c >= '0' && c <= '9' {
		d := int(c - '0')
		times := d + 1
		n.f[0] = State{exitSide: 1, exitDir: 1}
		n.f[1] = State{exitSide: 0, exitDir: 0}
		n.f[0].cnt[d] = times
		n.f[1].cnt[d] = times
	} else if c == '<' {
		n.f[0] = State{exitSide: 0, exitDir: 0}
		n.f[1] = State{exitSide: 0, exitDir: 0}
	} else { // '>'
		n.f[0] = State{exitSide: 1, exitDir: 1}
		n.f[1] = State{exitSide: 1, exitDir: 1}
	}
	return n
}

// combine A and B into C
func combine(A, B Node) Node {
	var C Node
	for entry := 0; entry < 2; entry++ {
		var cnt [10]int
		s := A.f[entry]
		for i := 0; i < 10; i++ {
			cnt[i] = s.cnt[i]
		}
		// first segment A
		if s.exitSide == 0 || s.exitSide == 1 && s.exitSide == 0 {
			// exit left or right directly
		}
		if s.exitSide == 0 {
			C.f[entry] = State{exitSide: 0, exitDir: s.exitDir, cnt: cnt}
		} else if s.exitSide == 1 {
			// enter B
			dir := s.exitDir
			t := B.f[0]
			if dir == 0 {
				// left
				t = B.f[1]
			}
			// accumulate B
			for i := 0; i < 10; i++ {
				cnt[i] += t.cnt[i]
			}
			if t.exitSide == 1 {
				C.f[entry] = State{exitSide: 1, exitDir: t.exitDir, cnt: cnt}
			} else {
				// bounce back to A once
				dir2 := t.exitDir
				u := A.f[0]
				if dir2 == 1 {
					u = A.f[1]
				}
				for i := 0; i < 10; i++ {
					cnt[i] += u.cnt[i]
				}
				C.f[entry] = State{exitSide: u.exitSide, exitDir: u.exitDir, cnt: cnt}
			}
		}
	}
	return C
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]Node, 2*size)
	// init
	for i := 0; i < n; i++ {
		seg[size+i] = makeNode(s[i])
	}
	// identity for others: zero node
	for i := size - 1; i > 0; i-- {
		seg[i] = combine(seg[2*i], seg[2*i+1])
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		// query [l,r)
		l += size
		r += size
		// two halves
		var leftRes *Node
		var rightRes *Node
		for l < r {
			if l&1 == 1 {
				if leftRes == nil {
					tmp := seg[l]
					leftRes = &tmp
				} else {
					tmp := combine(*leftRes, seg[l])
					leftRes = &tmp
				}
				l++
			}
			if r&1 == 1 {
				r--
				if rightRes == nil {
					tmp := seg[r]
					rightRes = &tmp
				} else {
					tmp := combine(seg[r], *rightRes)
					rightRes = &tmp
				}
			}
			l >>= 1
			r >>= 1
		}
		var res Node
		if leftRes == nil {
			if rightRes == nil {
				// empty
				res = Node{}
			} else {
				res = *rightRes
			}
		} else {
			if rightRes == nil {
				res = *leftRes
			} else {
				tmp := combine(*leftRes, *rightRes)
				res = tmp
			}
		}
		// start at entry 0
		for i := 0; i < 10; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", res.f[0].cnt[i]))
		}
		out.WriteByte('\n')
	}
}
