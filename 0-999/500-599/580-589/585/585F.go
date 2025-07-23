package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1000000007

// automaton node
type node struct {
	next [10]int
	fail int
	out  bool
}

var (
	nodes []node
	d     int
)

func newNode() int {
	n := len(nodes)
	var nd node
	for i := 0; i < 10; i++ {
		nd.next[i] = -1
	}
	nodes = append(nodes, nd)
	return n
}

func addPattern(s string) {
	v := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - '0')
		if nodes[v].next[c] == -1 {
			nodes[v].next[c] = newNode()
		}
		v = nodes[v].next[c]
	}
	nodes[v].out = true
}

func buildAC() {
	queue := make([]int, 0)
	for c := 0; c < 10; c++ {
		u := nodes[0].next[c]
		if u != -1 {
			nodes[u].fail = 0
			queue = append(queue, u)
		} else {
			nodes[0].next[c] = 0
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		f := nodes[v].fail
		for c := 0; c < 10; c++ {
			u := nodes[v].next[c]
			if u != -1 {
				nodes[u].fail = nodes[f].next[c]
				if nodes[nodes[u].fail].out {
					nodes[u].out = true
				}
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
}

func minusOne(s string) string {
	if s == "0" {
		return "0"
	}
	b := []byte(s)
	i := len(b) - 1
	for i >= 0 && b[i] == '0' {
		b[i] = '9'
		i--
	}
	if i >= 0 {
		b[i]--
	}
	j := 0
	for j < len(b)-1 && b[j] == '0' {
		j++
	}
	return string(b[j:])
}

func solve(num string) int {
	if len(num) < d {
		return 0
	}
	digits := make([]int, d)
	for i := 0; i < d; i++ {
		digits[i] = int(num[i] - '0')
	}
	stateCnt := len(nodes)
	dp := make([][][]int, d+1)
	vis := make([][][]bool, d+1)
	for i := 0; i <= d; i++ {
		dp[i] = make([][]int, stateCnt)
		vis[i] = make([][]bool, stateCnt)
		for j := 0; j < stateCnt; j++ {
			dp[i][j] = make([]int, 2)
			vis[i][j] = make([]bool, 2)
		}
	}
	var dfs func(pos, state, ok int, tight bool) int
	dfs = func(pos, state, ok int, tight bool) int {
		if pos == d {
			if ok == 1 {
				return 1
			}
			return 0
		}
		if !tight && vis[pos][state][ok] {
			return dp[pos][state][ok]
		}
		limit := 9
		if tight {
			limit = digits[pos]
		}
		start := 0
		if pos == 0 {
			start = 1
		}
		res := 0
		for dig := start; dig <= limit; dig++ {
			ns := nodes[state].next[dig]
			nOk := ok
			if nodes[ns].out {
				nOk = 1
			}
			t := tight && dig == limit
			val := dfs(pos+1, ns, nOk, t)
			res += val
			if res >= mod {
				res -= mod
			}
		}
		if !tight {
			vis[pos][state][ok] = true
			dp[pos][state][ok] = res
		}
		return res
	}
	return dfs(0, 0, 0, true)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s, x, y string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	fmt.Fscan(in, &x)
	fmt.Fscan(in, &y)
	d = len(x)
	L := d / 2
	nodes = nil
	newNode()
	if L > 0 && len(s) >= L {
		for i := 0; i+L <= len(s); i++ {
			addPattern(s[i : i+L])
		}
	}
	buildAC()
	resY := solve(y)
	xMinus := minusOne(x)
	resX := solve(xMinus)
	ans := resY - resX
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
