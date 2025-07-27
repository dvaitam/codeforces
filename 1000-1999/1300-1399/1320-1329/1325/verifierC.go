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

type edge struct{ u, v int }

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(n int, edges []edge) []int {
	inDeg := make([]int, n+1)
	for _, e := range edges {
		inDeg[e.u]++
		inDeg[e.v]++
	}
	node := 0
	for i := 1; i <= n; i++ {
		if inDeg[i] >= 3 {
			node = i
			break
		}
	}
	labels := make([]int, len(edges))
	for i := range labels {
		labels[i] = -1
	}
	cnt := 0
	if node != 0 {
		for i, e := range edges {
			if e.u == node || e.v == node {
				labels[i] = cnt
				cnt++
			}
		}
	}
	for i := range labels {
		if labels[i] == -1 {
			labels[i] = cnt
			cnt++
		}
	}
	return labels
}

func verifyCase(bin string, n int, edges []edge) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	outFields := strings.Fields(out)
	if len(outFields) != n-1 {
		return fmt.Errorf("expected %d numbers, got %d", n-1, len(outFields))
	}
	got := make([]int, n-1)
	for i, f := range outFields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		got[i] = v
	}
	expect := solveC(n, edges)
	for i := range expect {
		if got[i] != expect[i] {
			return fmt.Errorf("expected %v got %v", expect, got)
		}
	}
	return nil
}

func genTree(rng *rand.Rand, n int) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 2
		edges := genTree(rng, n)
		if err := verifyCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
