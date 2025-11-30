package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `2 3 0 1 1 0 0 0
1 4 1 0 0 1
3 2 0 1 0 0 1 1
2 2 1 1 1 0
5 3 1 0 0 0 1 1 0 1 0 1 1 0 1 1 1
1 1 1
2 4 0 1 0 1 1 1 1 0
2 2 0 1 0 1
1 3 1 0 1
2 3 0 1 1 0 0 1
5 6 0 1 0 1 0 0 0 1 0 0 0 0 1 1 1 0 1 0 0 1 1 1 1 1 0 1 0 1 1 1
3 3 1 1 1 0 1 0 1 1 1
1 5 0 0 0 0 1
2 2 1 0 1 1
4 4 1 1 0 1 1 1 1 1 0 0 0 0 1 1 1 0
5 6 1 0 0 0 0 1 1 0 0 0 0 1 1 0 0 0 0 1 1 0 0 0 0 1 0 0 0 1 0 0 0 0
2 2 1 1 1 1
3 6 1 1 1 1 1 1 1 1 1 1 1 1 1 0 1 1 0 1
2 2 1 1 1 1
1 3 1 1 1
3 5 0 0 1 1 1 0 0 0 0 1 1 0 0 1 1
2 5 1 0 1 1 1 0 0 1 0 0
1 5 1 1 1 1 0
1 5 1 0 1 1 0
5 4 1 0 0 1 0 1 1 1 0 0 0 1 1 0 0 1 1 1 0 1
4 4 1 1 1 0 1 1 1 0 1 1 0 1 0 1 1 0
3 5 1 0 0 0 0 0 1 0 1 0 0 0 0 0 0
2 3 0 0 1 1 1 1
3 5 0 1 1 0 0 1 1 1 1 1 1 1 1 0 0
2 3 1 1 0 1 1 1
5 6 0 0 0 1 0 0 1 0 1 1 1 1 0 0 0 0 0 0 1 1 1 1 1 0 1 0 1 0 0 0
5 6 0 0 1 1 1 1 1 1 0 0 0 1 0 1 1 0 0 0 1 1 0 0 1 0 1 0 1 1 1 1
3 5 0 0 1 0 1 0 1 0 1 0 1 0 0 0 0
2 6 1 0 1 1 0 0 0 0 1 0 0 0
2 2 1 0 0 1
3 3 0 0 1 1 1 1 1 1 1
2 3 1 0 0 1 0 1
2 4 0 0 0 0 0 0 1 1
5 3 0 0 0 1 0 1 0 1 0 0 0 0 1 1 0
3 3 0 1 0 0 0 0 0 1 0
3 4 0 0 0 0 0 0 0 0 0 0 0 0
1 5 0 1 1 1 0
1 2 1 0
1 1 1
2 2 0 1 1 0
2 3 0 0 0 1 0 0
2 4 1 0 1 1 1 0 1 1
2 3 1 1 0 1 1 0
1 6 1 1 1 0 1 0
3 5 0 0 1 1 1 0 1 0 0 1 1 1 0 1 0
1 6 0 1 1 0 1 0
2 6 0 1 1 1 1 1 0 1 1 1 1 0
2 2 1 0 0 0
2 5 0 0 1 1 0 1 1 1 0 0
2 2 1 1 1 1
1 2 0 0
1 2 1 1
4 5 1 1 1 0 1 1 0 0 0 1 1 0 1 1 1 1 1 0 1 0 0 0 0
3 5 0 1 0 0 1 1 0 0 0 0 1 1 0 1 0
1 5 0 1 0 1 1
3 6 0 0 0 1 1 0 1 1 0 0 0 0 0 1 1 0 0 1 0
2 4 0 1 0 0 0 0 0 0
3 6 0 1 1 1 0 0 0 1 0 0 1 0 0 0 1 0 1 1
2 4 0 0 0 1 0 0 1 1
5 5 0 0 0 0 1 1 0 0 1 0 1 0 0 0 1 1 1 1 0 1 0 0 1 1 0
4 6 1 1 1 1 1 0 1 1 0 0 0 1 0 0 1 1 1 0 1 1 0 1 1 1
2 5 1 1 0 1 0 1 0 0 0
1 5 0 1 0 1 1
2 2 1 1 1 1
1 5 0 1 0 0 1
5 6 1 1 1 0 0 1 0 0 0 0 1 0 1 0 0 0 0 1 1 0 0 1 0 0 1 1 1 1
5 6 0 1 1 0 0 0 0 0 0 0 0 1 1 1 1 0 1 1 0 1 1 0 0 0 0 0 1 0 0 1
5 5 1 0 0 1 1 1 0 0 0 1 0 0 1 0 1 0 1 1 0 0 1 1 1 0 1
1 5 0 1 0 1 1
5 6 0 0 0 0 0 0 1 0 1 1 0 0 1 1 1 1 1 1 1 0 1 1 1 1 1 1 1 1 1 0
1 5 1 0 1 0 0
3 3 0 0 0 0 0 0 1 1 1
1 2 1 1
4 5 0 1 0 1 1 0 1 0 1 0 1 1 1 1 0 1 1 1 0 1 1 1 0
2 5 1 1 0 1 0 0 1 1 0 0
1 5 0 0 0 0 1
4 6 0 0 0 1 1 1 0 1 1 0 1 0 1 0 0 0 0 0 1 0 0 0 1 1
5 6 1 0 0 1 1 1 1 1 1 0 0 0 0 1 1 1 0 1 1 1 0 0 1 0 0 0 0 0 0 1
5 6 1 1 0 0 1 1 0 0 0 0 1 1 0 0 0 1 1 0 1 1 1 0 0 0 1 1 1 1 1 1
5 6 1 1 0 1 0 1 0 1 0 1 0 0 0 0 0 0 1 1 1 0 1 1 0 1 1 1 1 1 0 0
2 6 0 0 1 0 0 0 0 0 0 1 0 1
2 3 0 1 0 1 0 1
3 3 1 0 0 0 0 0 1 0 0
5 6 0 0 0 1 0 0 0 1 0 0 0 1 1 0 1 0 1 0 0 1 1 1 0 0 1 0 0 1 0 1
5 5 0 0 0 1 0 1 1 0 1 1 1 1 1 1 1 0 0 0 1 1 0 1 0 0 0
2 6 1 1 0 0 0 0 0 1 1 0 0 0
3 2 1 0 0 0 0 1
5 3 1 1 1 0 0 0 1 1 1 1 1 1 1 1 1
2 2 0 0 0 0
4 5 0 0 1 1 1 1 0 1 0 1 0 0 0 1 1 1 1 1 1 0 1 1 1 0
4 5 0 0 0 1 0 1 0 0 0 0 0 0 1 1 1 0 1 0 1 0 1 0 0 0
1 5 1 0 0 0 1
2 5 0 0 1 1 1 1 1 0 1 0
1 4 0 0 0 1
1 5 1 0 1 0 1
4 5 0 1 1 0 1 0 1 1 0 0 0 1 1 1 0 1 1 1 0 0 0 1 1
1 3 1 1 1
1 6 1 1 0 1 1 1
5 5 0 0 0 1 1 1 1 1 0 0 1 1 1 1 1 0 0 0 1 1 0 1 1 1 0
5 5 0 0 1 0 0 0 0 1 1 1 0 1 0 1 1 1 0 0 1 0 1 0 0 1 0
4 6 1 1 0 1 1 1 0 0 0 1 0 0 1 0 1 1 1 0 1 1 1 1
1 6 1 1 1 1 1 0
5 6 0 0 0 0 1 0 0 1 0 0 1 0 1 1 0 1 0 0 0 0 1 0 1 0 1 1 1 0 1 0
5 5 0 0 1 1 1 1 0 0 1 0 1 0 1 0 0 0 1 1 0 1 1 0 0 1 1
4 6 0 1 0 1 0 0 0 1 0 0 0 1 0 1 0 0 0 1 0 1 1 1 1 1
2 3 1 0 1 1 1 0
5 5 1 0 1 0 1 1 1 0 0 1 1 1 0 1 1 1 1 0 1 1 0 1 0 1 1
4 6 0 0 1 0 0 0 1 0 1 0 1 1 0 1 0 1 1 0 1 1 1 0
5 6 1 1 0 1 0 0 0 1 1 0 0 0 0 1 1 1 1 1 0 0 1 1 1 1 1 0 1 0
3 3 1 1 1 0 0 0 0 0 0
1 4 1 1 0 1
5 6 0 0 1 1 1 0 0 0 1 0 1 0 1 0 1 0 1 0 0 0 0 1 1 1 0 1 1 1
1 6 0 0 1 1 0 1
2 2 0 1 1 0
2 2 1 1 1 1
1 6 0 0 1 0 1 0
2 4 0 0 1 1 0 0 0 0
2 4 0 0 0 0 0 1 1 1`

type edge struct {
	to  int
	id  int
	rev int
	used bool
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve36E(m int, pairs [][2]int) string {
	if m < 2 {
		return "-1"
	}
	maxv := 0
	for _, p := range pairs {
		if p[0] > maxv {
			maxv = p[0]
		}
		if p[1] > maxv {
			maxv = p[1]
		}
	}

	adj := make([][]*edge, maxv+1)
	for i, p := range pairs {
		id := i + 1
		u, v := p[0], p[1]
		eu := &edge{to: v, id: id, rev: len(adj[v])}
		ev := &edge{to: u, id: id, rev: len(adj[u])}
		adj[u] = append(adj[u], eu)
		adj[v] = append(adj[v], ev)
	}

	visitedV := make([]bool, maxv+1)
	var compsVerts [][]int
	var compsEdges [][]int
	for v := 1; v <= maxv; v++ {
		if !visitedV[v] && len(adj[v]) > 0 {
			stack := []int{v}
			visitedV[v] = true
			var compV []int
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				compV = append(compV, u)
				for _, e := range adj[u] {
					if !visitedV[e.to] {
						visitedV[e.to] = true
						stack = append(stack, e.to)
					}
				}
			}
			usedE := make([]bool, m+1)
			var compE []int
			for _, u := range compV {
				for _, e := range adj[u] {
					if !usedE[e.id] {
						usedE[e.id] = true
						compE = append(compE, e.id)
					}
				}
			}
			compsVerts = append(compsVerts, compV)
			compsEdges = append(compsEdges, compE)
		}
	}
	if len(compsVerts) > 2 {
		return "-1"
	}
	if len(compsVerts) == 2 {
		var trails [][]int
		for ci := 0; ci < 2; ci++ {
			compV := compsVerts[ci]
			compE := compsEdges[ci]
			localAdj := make([][]*edge, maxv+1)
			for _, id := range compE {
				u, v := pairs[id-1][0], pairs[id-1][1]
				eu := &edge{to: v, id: id, rev: len(localAdj[v])}
				ev := &edge{to: u, id: id, rev: len(localAdj[u])}
				localAdj[u] = append(localAdj[u], eu)
				localAdj[v] = append(localAdj[v], ev)
			}
			var localOdds []int
			for _, u := range compV {
				if len(localAdj[u])%2 == 1 {
					localOdds = append(localOdds, u)
				}
			}
			if len(localOdds) > 2 {
				return "-1"
			}
			startLocal := compV[0]
			if len(localOdds) == 2 {
				startLocal = localOdds[0]
			}
			ptr := make([]int, maxv+1)
			var path []int
			var dfs func(int)
			dfs = func(u int) {
				for ptr[u] < len(localAdj[u]) {
					e := localAdj[u][ptr[u]]
					ptr[u]++
					if e.used {
						continue
					}
					e.used = true
					localAdj[e.to][e.rev].used = true
					dfs(e.to)
					path = append(path, e.id)
				}
			}
			dfs(startLocal)
			if len(path) != len(compE) {
				return "-1"
			}
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			trails = append(trails, path)
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(len(trails[0])))
		sb.WriteByte('\n')
		for i, id := range trails[0] {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(id))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(len(trails[1])))
		sb.WriteByte('\n')
		for i, id := range trails[1] {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(id))
		}
		return sb.String()
	}

	odds := make([]int, 0, 4)
	for v := 1; v <= maxv; v++ {
		if len(adj[v])%2 == 1 {
			odds = append(odds, v)
		}
	}
	if len(odds) > 4 {
		return "-1"
	}
	dummyID := m + 1
	dummyAdded := false
	if len(odds) == 4 {
		u, v := odds[0], odds[1]
		eu := &edge{to: v, id: dummyID, rev: len(adj[v])}
		ev := &edge{to: u, id: dummyID, rev: len(adj[u])}
		adj[u] = append(adj[u], eu)
		adj[v] = append(adj[v], ev)
		dummyAdded = true
	}
	start := -1
	if dummyAdded {
		start = odds[2]
	} else if len(odds) == 2 {
		start = odds[0]
	} else {
		for v := 1; v <= maxv; v++ {
			if len(adj[v]) > 0 {
				start = v
				break
			}
		}
	}
	if start == -1 {
		return "-1"
	}
	ptr := make([]int, maxv+1)
	var path []int
	var dfs func(int)
	dfs = func(u int) {
		for ptr[u] < len(adj[u]) {
			e := adj[u][ptr[u]]
			ptr[u]++
			if e.used {
				continue
			}
			e.used = true
			adj[e.to][e.rev].used = true
			dfs(e.to)
			path = append(path, e.id)
		}
	}
	dfs(start)
	totalEdges := m
	if dummyAdded {
		totalEdges++
	}
	if len(path) != totalEdges {
		return "-1"
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	if dummyAdded {
		idx := -1
		for i, id := range path {
			if id == dummyID {
				idx = i
				break
			}
		}
		if idx == -1 || idx == 0 || idx == len(path)-1 {
			return "-1"
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(idx))
		sb.WriteByte('\n')
		for i := 0; i < idx; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(path[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(len(path) - idx - 1))
		sb.WriteByte('\n')
		for i := idx + 1; i < len(path); i++ {
			if i > idx+1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(path[i]))
		}
		return sb.String()
	}

	// split simple path into two non-empty parts
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(path[0]))
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(m - 1))
	sb.WriteByte('\n')
	for i := 1; i < len(path); i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(path[i]))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields)%2 != 0 {
			fmt.Fprintf(os.Stderr, "case %d: invalid field count\n", idx+1)
			os.Exit(1)
		}
		m := len(fields) / 2
		pairs := make([][2]int, m)
		pos := 0
		for i := 0; i < m; i++ {
			u, err1 := strconv.Atoi(fields[pos])
			v, err2 := strconv.Atoi(fields[pos+1])
			if err1 != nil || err2 != nil {
				fmt.Fprintf(os.Stderr, "case %d: parse error\n", idx+1)
				os.Exit(1)
			}
			pairs[i] = [2]int{u, v}
			pos += 2
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", m)
		for _, p := range pairs {
			fmt.Fprintf(&input, "%d %d\n", p[0], p[1])
		}
		want := strings.TrimSpace(solve36E(m, pairs))
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
