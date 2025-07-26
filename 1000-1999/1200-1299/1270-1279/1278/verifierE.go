package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func checkOutput(n int, edges [][2]int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	vals := make([]int, 0, 2*n)
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		vals = append(vals, v)
	}
	if len(vals) != 2*n {
		return fmt.Errorf("expected %d integers, got %d", 2*n, len(vals))
	}
	used := make([]bool, 2*n+1)
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		l := vals[2*i]
		r := vals[2*i+1]
		if l < 1 || l > 2*n || r < 1 || r > 2*n || l >= r {
			return fmt.Errorf("bad segment")
		}
		if used[l] || used[r] {
			return fmt.Errorf("duplicate endpoint")
		}
		used[l] = true
		used[r] = true
		L[i] = l
		R[i] = r
	}
	for v := 1; v <= 2*n; v++ {
		if !used[v] {
			return fmt.Errorf("missing endpoint %d", v)
		}
	}
	adj := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		adj[[2]int{e[0] - 1, e[1] - 1}] = true
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			inter := L[i] < R[j] && L[j] < R[i]
			contain := (L[i] < L[j] && R[j] < R[i]) || (L[j] < L[i] && R[i] < R[j])
			cross := inter && !contain
			key := [2]int{i, j}
			hasEdge := adj[key]
			if cross != hasEdge {
				return fmt.Errorf("pair %d %d mismatch", i+1, j+1)
			}
		}
	}
	return nil
}

func generateCase(r *rand.Rand) (string, [][2]int) {
	n := r.Intn(6) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, edges := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkOutput(len(edges)+1, edges, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
