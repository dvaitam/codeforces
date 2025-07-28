package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(n, m int, edges [][2]int) []int {
	g := make([][]int, n+1)
	rg := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		g[a] = append(g[a], b)
		rg[b] = append(rg[b], a)
	}
	reachable := make([]bool, n+1)
	stack := []int{1}
	reachable[1] = true
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if !reachable[to] {
				reachable[to] = true
				stack = append(stack, to)
			}
		}
	}
	order := make([]int, 0, n)
	used := make([]bool, n+1)
	var dfs1 func(int)
	dfs1 = func(v int) {
		used[v] = true
		for _, to := range g[v] {
			if reachable[to] && !used[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for v := 1; v <= n; v++ {
		if reachable[v] && !used[v] {
			dfs1(v)
		}
	}
	comp := make([]int, n+1)
	compCnt := 0
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for _, to := range rg[v] {
			if reachable[to] && comp[to] == 0 {
				dfs2(to, c)
			}
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == 0 {
			compCnt++
			dfs2(v, compCnt)
		}
	}
	compSize := make([]int, compCnt+1)
	for v := 1; v <= n; v++ {
		if reachable[v] {
			compSize[comp[v]]++
		}
	}
	cyc := make([]bool, compCnt+1)
	for c := 1; c <= compCnt; c++ {
		if compSize[c] > 1 {
			cyc[c] = true
		}
	}
	for v := 1; v <= n; v++ {
		if reachable[v] {
			cv := comp[v]
			for _, to := range g[v] {
				if reachable[to] && comp[to] == cv && v == to {
					cyc[cv] = true
				}
			}
		}
	}
	adjC := make([][]int, compCnt+1)
	for v := 1; v <= n; v++ {
		if !reachable[v] {
			continue
		}
		cv := comp[v]
		for _, to := range g[v] {
			if !reachable[to] {
				continue
			}
			ct := comp[to]
			if cv != ct {
				adjC[cv] = append(adjC[cv], ct)
			}
		}
	}
	startComp := comp[1]
	reachComp := make([]bool, compCnt+1)
	queue := []int{startComp}
	reachComp[startComp] = true
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, to := range adjC[c] {
			if !reachComp[to] {
				reachComp[to] = true
				queue = append(queue, to)
			}
		}
	}
	inf := make([]bool, compCnt+1)
	q := make([]int, 0)
	for c := 1; c <= compCnt; c++ {
		if reachComp[c] && cyc[c] {
			inf[c] = true
			q = append(q, c)
		}
	}
	for head := 0; head < len(q); head++ {
		c := q[head]
		for _, to := range adjC[c] {
			if reachComp[to] && !inf[to] {
				inf[to] = true
				q = append(q, to)
			}
		}
	}
	dp := make([]int, compCnt+1)
	indeg := make([]int, compCnt+1)
	for c := 1; c <= compCnt; c++ {
		if !reachComp[c] || inf[c] {
			continue
		}
		for _, to := range adjC[c] {
			if reachComp[to] && !inf[to] {
				indeg[to]++
			}
		}
	}
	queue = queue[:0]
	for c := 1; c <= compCnt; c++ {
		if !reachComp[c] || inf[c] {
			continue
		}
		if indeg[c] == 0 {
			queue = append(queue, c)
		}
	}
	dp[startComp] = 1
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, to := range adjC[c] {
			if !reachComp[to] || inf[to] {
				continue
			}
			val := dp[to] + dp[c]
			if val > 2 {
				val = 2
			}
			if val > dp[to] {
				dp[to] = val
			}
			indeg[to]--
			if indeg[to] == 0 {
				queue = append(queue, to)
			}
		}
	}
	ans := make([]int, n+1)
	for v := 1; v <= n; v++ {
		if !reachable[v] {
			ans[v] = 0
		} else if inf[comp[v]] {
			ans[v] = -1
		} else {
			ans[v] = dp[comp[v]]
		}
	}
	return ans[1:]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesG.txt")
	if err != nil {
		fmt.Println("could not read testcasesG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{a, b}
		}
		ans := solveCase(n, m, edges)
		strs := make([]string, len(ans))
		for j, v := range ans {
			strs[j] = strconv.Itoa(v)
		}
		expected[i] = strings.Join(strs, " ")
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanLines)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := strings.TrimSpace(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
