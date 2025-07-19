package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ st, ed int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var str string
	if _, err := fmt.Fscan(reader, &str); err != nil {
		return
	}
	n := len(str)
	// 1-based indexing for s
	s := make([]byte, n+1)
	for i := 1; i <= n; i++ {
		s[i] = str[i-1]
	}
	// next and pre links
	pre := make([]int, n+2)
	next := make([]int, n+2)
	vis := make([]bool, n+2)
	// stacks for two types: 0 for L, 1 for R
	stacks := [2][]int{}
	// build initial chains
	for i := 1; i <= n; i++ {
		p := 0
		if s[i] == 'R' {
			p = 1
		}
		opp := p ^ 1
		if len(stacks[opp]) > 0 {
			j := stacks[opp][len(stacks[opp])-1]
			stacks[opp] = stacks[opp][:len(stacks[opp])-1]
			vis[i] = true
			pre[i] = j
			next[j] = i
		} else {
			stacks[p] = append(stacks[p], i)
		}
	}
	// sequences of (start, end)
	seq := make([]pair, 0, n)
	// collect heads
	for i := 1; i <= n; i++ {
		if vis[i] {
			continue
		}
		// i is head
		j := i
		for next[j] != 0 {
			j = next[j]
		}
		seq = append(seq, pair{i, j})
		// try merging greedily
		for {
			var ok bool
			seq, ok = mergeSeq(seq, s, pre, next)
			if !ok {
				break
			}
		}
	}
	cnt := len(seq) - 1
	fmt.Fprintln(writer, cnt)
	// print sequences
	for _, t := range seq {
		for j := t.st; ; j = next[j] {
			fmt.Fprintf(writer, "%d", j)
			if j == t.ed {
				break
			}
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, " ")
	}
}

// mergeSeq tries to merge last two sequences, returns new seq and true if merged
func mergeSeq(seq []pair, s []byte, pre, next []int) ([]pair, bool) {
	if len(seq) < 2 {
		return seq, false
	}
	m := len(seq)
	st, ed := seq[m-1].st, seq[m-1].ed
	st2, ed2 := seq[m-2].st, seq[m-2].ed
	add := func(x, y int) {
		pre[y] = x
		next[x] = y
	}
	// case 1
	if s[st] != s[ed2] {
		add(ed2, st)
		// merge into (st2, ed)
		seq = seq[:m-2]
		seq = append(seq, pair{st2, ed})
		return seq, true
	}
	// case 2
	if s[st2] != s[ed] {
		add(ed, st2)
		seq = seq[:m-2]
		seq = append(seq, pair{st, ed2})
		return seq, true
	}
	// case 3
	if s[st] != s[ed] {
		if ed < ed2 {
			z := pre[ed2]
			add(ed, ed2)
			add(ed2, st2)
			next[z] = 0
			seq = seq[:m-2]
			seq = append(seq, pair{st, z})
		} else {
			z := pre[ed]
			add(ed2, ed)
			add(ed, st)
			next[z] = 0
			seq = seq[:m-2]
			seq = append(seq, pair{st2, z})
		}
		return seq, true
	}
	return seq, false
}
