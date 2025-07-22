package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct {
	u, v, w int
}

func solve(n int, edges []edge, k int) (int, bool) {
	m := len(edges)
	bestCost := int(^uint(0) >> 1)
	found := false
	for mask := 0; mask < (1 << m); mask++ {
		if bits.OnesCount(uint(mask)) != n-1 {
			continue
		}
		parent := make([]int, n+1)
		for i := range parent {
			parent[i] = i
		}
		var find func(int) int
		find = func(x int) int {
			if parent[x] != x {
				parent[x] = find(parent[x])
			}
			return parent[x]
		}
		union := func(a, b int) bool {
			fa, fb := find(a), find(b)
			if fa == fb {
				return false
			}
			parent[fb] = fa
			return true
		}
		cost := 0
		capCnt := 0
		edgesUsed := 0
		valid := true
		for i := 0; i < m; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			e := edges[i]
			if !union(e.u, e.v) {
				valid = false
				break
			}
			if e.u == 1 || e.v == 1 {
				capCnt++
			}
			cost += e.w
			edgesUsed++
		}
		if !valid || edgesUsed != n-1 || capCnt != k {
			continue
		}
		root := find(1)
		for i := 2; i <= n; i++ {
			if find(i) != root {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		if cost < bestCost {
			bestCost = cost
			found = true
		}
	}
	return bestCost, found
}

func checkOutput(n int, edges []edge, k int, bestCost int, hasSol bool, out string) error {
	out = strings.TrimSpace(out)
	if !hasSol {
		if out != "-1" {
			return fmt.Errorf("expected -1, got %q", out)
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected two lines in output")
	}
	firstVal, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || firstVal != n-1 {
		return fmt.Errorf("first line should be %d", n-1)
	}
	ids := strings.Fields(lines[1])
	if len(ids) != n-1 {
		return fmt.Errorf("expected %d edge indices", n-1)
	}
	used := make([]bool, len(edges))
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) bool {
		fa, fb := find(a), find(b)
		if fa == fb {
			return false
		}
		parent[fb] = fa
		return true
	}
	cost := 0
	capCnt := 0
	for _, s := range ids {
		idx, err := strconv.Atoi(s)
		if err != nil || idx < 1 || idx > len(edges) {
			return fmt.Errorf("invalid edge index %q", s)
		}
		if used[idx-1] {
			return fmt.Errorf("duplicate edge index")
		}
		used[idx-1] = true
		e := edges[idx-1]
		if !union(e.u, e.v) {
			return fmt.Errorf("edges do not form a tree")
		}
		if e.u == 1 || e.v == 1 {
			capCnt++
		}
		cost += e.w
	}
	if capCnt != k {
		return fmt.Errorf("capital edge count mismatch")
	}
	root := find(1)
	for i := 2; i <= n; i++ {
		if find(i) != root {
			return fmt.Errorf("graph not connected")
		}
	}
	if cost != bestCost {
		return fmt.Errorf("cost is not minimal")
	}
	return nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			w, _ := strconv.Atoi(scan.Text())
			edges[i] = edge{u, v, w}
		}
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		for _, e := range edges {
			input += fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w)
		}
		bestCost, ok := solve(n, edges, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if err := checkOutput(n, edges, k, bestCost, ok, out); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
