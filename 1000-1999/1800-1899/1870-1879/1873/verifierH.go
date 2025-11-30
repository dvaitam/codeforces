package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
5 2 4
1 2
1 3
3 4
1 5
3 5
3 3 1
1 2
1 3
3 2
3 1 3
1 2
2 3
3 1
7 3 5
1 2
1 3
3 4
2 5
3 6
1 7
5 6
3 3 1
1 2
2 3
1 3
4 1 3
1 2
2 3
3 4
4 2
5 2 4
1 2
1 3
3 4
1 5
3 5
6 5 4
1 2
1 3
2 4
4 5
1 6
6 3
8 8 5
1 2
2 3
1 4
4 5
3 6
2 7
5 8
2 8
3 1 2
1 2
1 3
2 3
5 2 4
1 2
2 3
2 4
4 5
3 1
4 4 3
1 2
1 3
3 4
1 4
8 7 2
1 2
2 3
1 4
2 5
1 6
2 7
4 8
2 4
7 4 2
1 2
2 3
3 4
3 5
4 6
1 7
3 1
4 1 3
1 2
2 3
3 4
3 1
8 5 6
1 2
2 3
1 4
3 5
2 6
5 7
7 8
7 4
4 4 3
1 2
1 3
2 4
4 3
4 3 4
1 2
2 3
1 4
2 4
4 3 2
1 2
1 3
2 4
1 4
8 4 8
1 2
2 3
3 4
3 5
1 6
6 7
4 8
8 7
8 2 3
1 2
1 3
1 4
2 5
5 6
4 7
7 8
3 8
8 6 3
1 2
1 3
1 4
1 5
5 6
6 7
2 8
7 4
4 1 3
1 2
2 3
3 4
4 2
7 2 5
1 2
2 3
1 4
1 5
2 6
2 7
2 4
7 6 1
1 2
2 3
3 4
4 5
1 6
5 7
1 7
3 3 2
1 2
1 3
2 3
4 3 2
1 2
1 3
2 4
1 4
6 3 1
1 2
2 3
1 4
2 5
3 6
2 4
4 1 4
1 2
1 3
3 4
2 3
5 1 3
1 2
2 3
3 4
4 5
4 1
6 3 5
1 2
1 3
1 4
2 5
1 6
2 3
4 4 3
1 2
1 3
3 4
4 1
5 1 3
1 2
1 3
1 4
3 5
5 4
5 5 2
1 2
1 3
1 4
2 5
3 4
7 6 2
1 2
2 3
1 4
3 5
1 6
1 7
2 5
6 2 4
1 2
2 3
3 4
4 5
5 6
3 6
4 2 3
1 2
1 3
2 4
4 3
7 2 6
1 2
1 3
2 4
2 5
2 6
3 7
4 1
5 3 5
1 2
1 3
1 4
3 5
2 3
4 1 3
1 2
1 3
2 4
3 2
4 1 3
1 2
1 3
2 4
1 4
3 2 3
1 2
1 3
3 2
4 3 2
1 2
2 3
3 4
1 4
6 5 1
1 2
1 3
2 4
3 5
1 6
4 1
8 2 8
1 2
1 3
2 4
2 5
2 6
2 7
6 8
8 5
3 3 1
1 2
1 3
3 2
7 5 2
1 2
2 3
1 4
4 5
3 6
6 7
1 6
4 4 3
1 2
2 3
2 4
4 1
7 2 3
1 2
2 3
1 4
3 5
4 6
1 7
7 5
6 3 4
1 2
1 3
1 4
1 5
2 6
6 5
5 3 2
1 2
1 3
3 4
3 5
2 4
6 4 1
1 2
1 3
2 4
4 5
4 6
3 6
4 4 3
1 2
2 3
1 4
3 1
5 3 4
1 2
1 3
3 4
1 5
5 3
4 3 4
1 2
2 3
3 4
1 4
6 5 2
1 2
2 3
1 4
2 5
2 6
6 3
6 6 2
1 2
2 3
3 4
4 5
1 6
2 6
4 1 2
1 2
1 3
2 4
3 4
6 2 5
1 2
1 3
1 4
2 5
3 6
5 1
5 2 3
1 2
1 3
1 4
4 5
5 2
6 3 1
1 2
2 3
3 4
3 5
2 6
6 5
7 6 7
1 2
1 3
2 4
2 5
4 6
4 7
4 3
3 1 2
1 2
2 3
1 3
8 2 8
1 2
1 3
1 4
2 5
2 6
5 7
1 8
5 3
8 5 7
1 2
1 3
1 4
3 5
5 6
5 7
2 8
1 5
6 3 6
1 2
2 3
3 4
2 5
5 6
4 6
8 5 1
1 2
1 3
2 4
4 5
1 6
3 7
2 8
7 6
4 4 1
1 2
2 3
2 4
1 3
8 2 4
1 2
1 3
2 4
2 5
2 6
4 7
2 8
3 4
6 4 6
1 2
1 3
2 4
1 5
2 6
1 6
3 1 2
1 2
2 3
3 1
3 1 2
1 2
1 3
3 2
3 2 3
1 2
2 3
1 3
4 4 3
1 2
2 3
1 4
4 2
5 3 4
1 2
2 3
1 4
4 5
1 5
5 3 5
1 2
2 3
3 4
3 5
1 5
8 2 1
1 2
1 3
2 4
4 5
4 6
3 7
4 8
8 3
6 2 1
1 2
1 3
3 4
2 5
3 6
3 5
3 3 1
1 2
1 3
3 2
4 2 4
1 2
1 3
3 4
2 4
4 2 3
1 2
2 3
1 4
4 3
4 2 1
1 2
1 3
1 4
2 4
5 1 3
1 2
1 3
1 4
2 5
3 5
4 4 3
1 2
1 3
3 4
3 2
3 2 1
1 2
1 3
2 3
5 2 5
1 2
1 3
1 4
1 5
4 5
6 1 4
1 2
2 3
3 4
2 5
5 6
1 6
4 4 3
1 2
2 3
3 4
1 3
8 6 7
1 2
1 3
2 4
2 5
4 6
6 7
4 8
4 1
6 2 4
1 2
1 3
2 4
3 5
4 6
5 2
8 7 2
1 2
1 3
1 4
3 5
3 6
2 7
5 8
3 2
3 2 1
1 2
1 3
3 2
7 7 2
1 2
1 3
3 4
2 5
1 6
4 7
2 6
4 1 3
1 2
2 3
3 4
2 4
6 6 3
1 2
2 3
1 4
1 5
1 6
5 4
6 2 4
1 2
1 3
2 4
4 5
1 6
3 4
5 1 4
1 2
1 3
3 4
2 5
5 1
3 3 2
1 2
1 3
2 3
5 2 3
1 2
2 3
1 4
3 5
5 4
4 3 2
1 2
2 3
1 4
4 2`

func expected(n, a, b int, edges [][2]int) string {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	queue := make([]int, 0, n)
	removed := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		removed[v] = true
		for _, to := range adj[v] {
			if removed[to] {
				continue
			}
			deg[to]--
			if deg[to] == 1 {
				queue = append(queue, to)
			}
		}
	}
	onCycle := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if !removed[i] {
			onCycle[i] = true
		}
	}
	distA := make([]int, n+1)
	for i := range distA {
		distA[i] = -1
	}
	qa := []int{a}
	distA[a] = 0
	for h := 0; h < len(qa); h++ {
		v := qa[h]
		for _, to := range adj[v] {
			if distA[to] == -1 {
				distA[to] = distA[v] + 1
				qa = append(qa, to)
			}
		}
	}
	distB := make([]int, n+1)
	for i := range distB {
		distB[i] = -1
	}
	qb := []int{b}
	distB[b] = 0
	for h := 0; h < len(qb); h++ {
		v := qb[h]
		for _, to := range adj[v] {
			if distB[to] == -1 {
				distB[to] = distB[v] + 1
				qb = append(qb, to)
			}
		}
	}
	if distB[b] >= distA[b] {
		return "NO\n"
	}
	visited := make([]bool, n+1)
	q := []int{b}
	visited[b] = true
	escape := false
	for head := 0; head < len(q); head++ {
		v := q[head]
		if onCycle[v] {
			escape = true
			break
		}
		for _, to := range adj[v] {
			if !visited[to] && distB[to] < distA[to] {
				visited[to] = true
				q = append(q, to)
			}
		}
	}
	if escape {
		return "YES\n"
	}
	return "NO\n"
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	pos := 0
	t, err := strconv.Atoi(tokens[pos])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	pos++
	var inputs []string
	var expects []string
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if pos+2 >= len(tokens) {
			fmt.Printf("case %d missing header\n", caseIdx)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		a, errA := strconv.Atoi(tokens[pos+1])
		b, errB := strconv.Atoi(tokens[pos+2])
		if errN != nil || errA != nil || errB != nil {
			fmt.Printf("case %d invalid header\n", caseIdx)
			os.Exit(1)
		}
		pos += 3
		if pos+2*n > len(tokens) {
			fmt.Printf("case %d missing edges\n", caseIdx)
			os.Exit(1)
		}
		edges := make([][2]int, n)
		for i := 0; i < n; i++ {
			u, errU := strconv.Atoi(tokens[pos])
			v, errV := strconv.Atoi(tokens[pos+1])
			if errU != nil || errV != nil {
				fmt.Printf("case %d invalid edge\n", caseIdx)
				os.Exit(1)
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", n, a, b)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		inputs = append(inputs, sb.String())
		expects = append(expects, expected(n, a, b, edges))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		if err := runCase(exe, input, expects[idx]); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
