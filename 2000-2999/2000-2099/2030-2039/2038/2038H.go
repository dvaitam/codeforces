package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -1 << 60

type state struct {
	turn   uint8
	prev   int8
	counts [50]uint8
}

type result struct {
	score  int64
	choice int8
}

var (
	n, m int
	p    []int
	a    [][]int64
	memo map[state]result
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n, &m)
	p = make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &p[i])
		p[i]--
	}

	a = make([][]int64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}

	memo = make(map[state]result)
	counts := make([]uint8, n)
	best := solve(0, -1, counts)
	if best == negInf {
		fmt.Fprintln(out, -1)
		return
	}

	ans := make([]int, m)
	prev := -1
	for turn := 0; turn < m; turn++ {
		st := makeState(turn, prev, counts)
		res := memo[st]
		choice := int(res.choice)
		ans[turn] = choice + 1
		counts[choice]++
		prev = leader(counts)
	}

	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}

func solve(turn int, prev int, counts []uint8) int64 {
	st := makeState(turn, prev, counts)
	if val, ok := memo[st]; ok {
		return val.score
	}
	if turn == m {
		memo[st] = result{0, -1}
		return 0
	}

	target := p[turn]
	best := negInf
	var bestChoice int8 = -1

	for c := 0; c < n; c++ {
		if prev != -1 && c == prev {
			continue
		}
		counts[c]++
		lead := leader(counts)
		if lead == target {
			val := solve(turn+1, lead, counts)
			if val != negInf {
				val += a[c][turn]
				if val > best {
					best = val
					bestChoice = int8(c)
				}
			}
		}
		counts[c]--
	}

	memo[st] = result{best, bestChoice}
	return best
}

func leader(counts []uint8) int {
	lead := 0
	mx := counts[0]
	for i := 1; i < n; i++ {
		if counts[i] > mx || (counts[i] == mx && i < lead) {
			mx = counts[i]
			lead = i
		}
	}
	return lead
}

func makeState(turn int, prev int, counts []uint8) state {
	var st state
	st.turn = uint8(turn)
	st.prev = int8(prev)
	for i := 0; i < n && i < 50; i++ {
		st.counts[i] = counts[i]
	}
	return st
}
