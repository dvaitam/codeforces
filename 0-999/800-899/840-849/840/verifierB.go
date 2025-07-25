package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Edge struct{ u, v int }

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solutionExists(d []int) bool {
	sum := 0
	hasNeg := false
	for _, v := range d {
		if v == -1 {
			hasNeg = true
		} else {
			sum += v
		}
	}
	if !hasNeg && sum%2 == 1 {
		return false
	}
	return true
}

func checkSolution(n, m int, d []int, edges []Edge, out string, exists bool) error {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	if len(lines) == 0 || lines[0] == "" {
		return fmt.Errorf("empty output")
	}
	if lines[0] == "-1" {
		if exists {
			return fmt.Errorf("declared no solution but one exists")
		}
		if len(lines) > 1 && strings.TrimSpace(strings.Join(lines[1:], "")) != "" {
			return fmt.Errorf("extraneous output after -1")
		}
		return nil
	}
	k, err := strconv.Atoi(strings.Fields(lines[0])[0])
	if err != nil {
		return fmt.Errorf("invalid first line")
	}
	if k < 0 || k > m {
		return fmt.Errorf("invalid k")
	}
	if len(lines)-1 < k {
		return fmt.Errorf("not enough edge indices")
	}
	used := make(map[int]bool)
	deg := make([]int, n+1)
	for i := 0; i < k; i++ {
		idxStr := strings.TrimSpace(lines[i+1])
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return fmt.Errorf("invalid edge index")
		}
		if idx < 1 || idx > m {
			return fmt.Errorf("edge index out of range")
		}
		if used[idx] {
			return fmt.Errorf("duplicate edge index")
		}
		used[idx] = true
		e := edges[idx-1]
		deg[e.u]++
		deg[e.v]++
	}
	for i := 1; i <= n; i++ {
		if d[i] == -1 {
			continue
		}
		if deg[i]%2 != d[i] {
			return fmt.Errorf("vertex %d parity mismatch", i)
		}
	}
	if !exists {
		return fmt.Errorf("solution claimed but none should exist")
	}
	return nil
}

func genCase(rng *rand.Rand) (string, int, int, []int, []Edge, bool) {
	n := rng.Intn(5) + 2
	m := n - 1 + rng.Intn(3) // up to n+1 edges
	d := make([]int, n+1)
	hasNeg := false
	for i := 1; i <= n; i++ {
		t := rng.Intn(3) - 1 // -1,0,1
		d[i] = t
		if t == -1 {
			hasNeg = true
		}
	}
	edges := make([]Edge, 0, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, Edge{p, i})
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, Edge{u, v})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(d[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	exists := solutionExists(d[1:]) || hasNeg
	return sb.String(), n, m, d, edges, exists
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, m, d, edges, exists := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkSolution(n, m, d, edges, out, exists); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
