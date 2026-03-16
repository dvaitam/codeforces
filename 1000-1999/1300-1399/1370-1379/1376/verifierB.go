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

const testcasesRaw = `2 0
3 3 2 3 1 3 1 2
4 0
4 2 3 4 1 3
6 13 3 6 2 4 1 4 4 6 2 5 1 6 3 5 2 6 5 6 1 2 2 3 1 3 3 4
6 13 1 4 5 6 1 6 3 4 2 4 3 5 1 3 2 3 4 5 4 6 3 6 2 5 1 5
5 4 4 5 2 5 3 4 1 5
3 3 1 2 1 3 2 3
5 1 4 5
4 2 1 4 1 3
4 5 1 4 1 3 2 4 2 3 1 2
5 6 2 4 3 4 1 5 1 2 3 5 2 3
5 8 2 5 1 5 1 3 1 2 2 3 3 5 1 4 3 4
5 7 3 4 3 5 2 4 1 3 2 3 1 4 4 5
3 0
5 1 1 5
5 4 1 4 4 5 3 4 2 4
1 0
3 1 2 3
2 1 1 2
5 0
3 3 1 2 2 3 1 3
2 1 1 2
1 0
6 7 2 3 1 6 1 4 4 5 2 5 5 6 1 3
6 4 4 6 1 3 3 4 3 5
2 0
3 2 1 3 2 3
3 3 1 3 1 2 2 3
5 2 2 4 1 5
6 12 1 4 4 5 3 6 4 6 2 5 5 6 1 6 3 5 1 3 2 3 1 5 2 4
3 0
3 2 1 2 1 3
1 0
2 1 1 2
5 10 2 3 3 4 3 5 1 5 1 2 4 5 1 3 2 5 2 4 1 4
5 2 1 3 1 4
5 1 1 2
3 1 2 3
5 5 4 5 1 4 2 3 3 4 2 5
1 0
1 0
5 1 1 5
1 0
1 0
1 0
1 0
2 0
5 6 4 5 1 2 3 5 2 3 2 4 3 4
4 2 1 2 1 3
1 0
6 14 2 6 1 3 3 5 3 6 2 4 1 4 1 5 2 5 5 6 3 4 1 6 4 5 4 6 1 2
1 0
2 1 1 2
6 10 3 4 2 5 1 3 5 6 2 4 1 2 4 5 3 6 3 5 4 6
2 1 1 2
3 1 2 3
3 1 1 2
3 1 1 3
2 0
2 1 1 2
1 0
5 2 2 5 1 5
1 0
1 0
2 0
2 1 1 2
3 1 2 3
2 0
5 4 2 3 3 5 4 5 2 5
1 0
4 0
4 6 2 3 1 2 2 4 3 4 1 3 1 4
6 5 3 5 4 6 1 6 1 5 4 5
5 10 2 5 2 4 3 5 1 4 2 3 4 5 1 2 3 4 1 5 1 3
6 1 1 3
1 0
2 0
3 2 1 2 2 3
2 1 1 2
5 7 2 4 2 3 1 2 1 5 1 4 3 4 3 5
6 4 4 5 1 5 1 3 4 6
5 7 1 4 4 5 1 2 2 3 2 4 3 5 2 5
2 1 1 2
2 0
4 2 1 4 1 3
2 1 1 2
6 9 2 4 4 5 3 5 1 4 1 3 1 6 2 3 3 6 4 6
2 0
3 3 1 2 2 3 1 3
2 1 1 2
4 1 2 4
3 0
5 9 1 5 3 4 1 4 4 5 1 3 2 5 2 4 2 3 3 5
3 2 1 2 2 3
1 0
4 2 1 2 2 4
6 9 2 3 2 6 4 5 3 6 1 4 1 2 1 6 1 5 4 6
6 2 4 6 5 6
4 3 2 4 1 2 2 3`

type graph struct {
	n     int
	edges [][2]int
	adj   [][]bool
}

func parseGraph(line string) (graph, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return graph{}, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return graph{}, err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return graph{}, err
	}
	if len(fields) != 2+2*m {
		return graph{}, fmt.Errorf("expected %d edge numbers, got %d", 2*m, len(fields)-2)
	}
	g := graph{n: n, edges: make([][2]int, m), adj: make([][]bool, n+1)}
	for i := range g.adj {
		g.adj[i] = make([]bool, n+1)
	}
	idx := 2
	for i := 0; i < m; i++ {
		a, _ := strconv.Atoi(fields[idx])
		b, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		g.edges[i] = [2]int{a, b}
		g.adj[a][b] = true
		g.adj[b][a] = true
	}
	return g, nil
}

func graphInput(g graph) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", g.n, len(g.edges)))
	for _, e := range g.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkOutput(g graph, out string) error {
	fields := strings.Fields(out)
	if len(fields) != 1+g.n {
		return fmt.Errorf("expected %d numbers, got %d", 1+g.n, len(fields))
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	choose := make([]int, g.n+1)
	cnt := 0
	for i := 1; i <= g.n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("bad value for vertex %d", i)
		}
		if v != 0 && v != 1 {
			return fmt.Errorf("vertex %d not 0/1", i)
		}
		choose[i] = v
		if v == 1 {
			cnt++
		}
	}
	if cnt != k {
		return fmt.Errorf("claimed size %d but got %d", k, cnt)
	}
	// independence check
	for _, e := range g.edges {
		if choose[e[0]] == 1 && choose[e[1]] == 1 {
			return fmt.Errorf("edge %d-%d violates independence", e[0], e[1])
		}
	}
	// maximality check
	for i := 1; i <= g.n; i++ {
		if choose[i] == 1 {
			continue
		}
		canAdd := true
		for j := 1; j <= g.n; j++ {
			if choose[j] == 1 && g.adj[i][j] {
				canAdd = false
				break
			}
		}
		if canAdd {
			return fmt.Errorf("vertex %d could be added", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		g, err := parseGraph(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test case %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := graphInput(g)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if err := checkOutput(g, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", idx, err, input, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
