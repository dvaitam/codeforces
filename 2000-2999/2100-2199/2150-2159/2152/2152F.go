package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func nextInt() int {
	sign := 1
	val := 0
	c, err := reader.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = reader.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = reader.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = reader.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCase() {
	n := nextInt()
	z := nextInt()
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = nextInt()
	}

	limit := make([]int, n)
	ptr := 0
	for i := 0; i < n; i++ {
		if ptr < i+1 {
			ptr = i + 1
		}
		threshold := int64(x[i]) + int64(z)
		for ptr < n && int64(x[ptr]) <= threshold {
			ptr++
		}
		limit[i] = ptr
	}

	type pair struct {
		p, q int
	}
	states := make([]pair, 0, n)
	nextState := make([]int, 0, n)
	idMap := make(map[int64]int, n*2+2)
	base := int64(n + 1)

	var getID func(p, q int) int
	getID = func(p, q int) int {
		key := int64(p)*base + int64(q)
		if v, ok := idMap[key]; ok {
			return v
		}
		id := len(states)
		states = append(states, pair{p: p, q: q})
		idMap[key] = id
		return id
	}

	for l := 0; l+1 < n; l++ {
		getID(l, l+1)
	}

	for idx := 0; idx < len(states); idx++ {
		p := states[idx].p
		q := states[idx].q
		nxt := max(limit[p], q+1)
		if nxt < n {
			nid := getID(q, nxt)
			nextState = append(nextState, nid)
		} else {
			nextState = append(nextState, -1)
		}
	}

	const LOG = 20
	m := len(states)
	up := make([][]int, LOG)
	up[0] = make([]int, m)
	copy(up[0], nextState)
	for k := 1; k < LOG; k++ {
		up[k] = make([]int, m)
		for i := 0; i < m; i++ {
			prev := up[k-1][i]
			if prev == -1 {
				up[k][i] = -1
			} else {
				up[k][i] = up[k-1][prev]
			}
		}
	}

	qQueries := nextInt()
	for ; qQueries > 0; qQueries-- {
		l := nextInt() - 1
		r := nextInt() - 1
		if l == r {
			fmt.Fprintln(writer, 1)
			continue
		}
		if l+1 > r {
			fmt.Fprintln(writer, 1)
			continue
		}
		key := int64(l)*base + int64(l+1)
		state, ok := idMap[key]
		if !ok {
			fmt.Fprintln(writer, 1)
			continue
		}
		if states[state].q > r {
			fmt.Fprintln(writer, 1)
			continue
		}
		ans := 2
		for k := LOG - 1; k >= 0; k-- {
			nxt := up[k][state]
			if nxt != -1 && states[nxt].q <= r {
				ans += 1 << k
				state = nxt
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

func main() {
	defer writer.Flush()
	t := nextInt()
	for ; t > 0; t-- {
		solveCase()
	}
}
