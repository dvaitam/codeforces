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

func checkPaths(n, k, a, b, c, d int, v, u []int) error {
	if len(v) != n || len(u) != n {
		return fmt.Errorf("expected %d numbers per line", n)
	}
	if v[0] != a || v[n-1] != b || u[0] != c || u[n-1] != d {
		return fmt.Errorf("endpoints mismatch")
	}
	used1 := make([]bool, n+1)
	used2 := make([]bool, n+1)
	for i := 0; i < n; i++ {
		if v[i] < 1 || v[i] > n || used1[v[i]] {
			return fmt.Errorf("line1 not a permutation")
		}
		used1[v[i]] = true
		if u[i] < 1 || u[i] > n || used2[u[i]] {
			return fmt.Errorf("line2 not a permutation")
		}
		used2[u[i]] = true
	}
	edges := make(map[[2]int]struct{})
	addEdge := func(x, y int) {
		if x > y {
			x, y = y, x
		}
		edges[[2]int{x, y}] = struct{}{}
	}
	for i := 0; i < n-1; i++ {
		addEdge(v[i], v[i+1])
		addEdge(u[i], u[i+1])
	}
	if _, ok := edges[[2]int{a, b}]; ok {
		return fmt.Errorf("edge between a and b present")
	}
	if _, ok := edges[[2]int{c, d}]; ok {
		return fmt.Errorf("edge between c and d present")
	}
	if len(edges) > k {
		return fmt.Errorf("too many edges: %d > %d", len(edges), k)
	}
	return nil
}

func runCase(exe string, input string, n, k, a, b, c, d int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	if output == "-1" {
		if n <= 4 || k < n+1 {
			return nil
		}
		return fmt.Errorf("reported no solution, but one exists")
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected two lines of output")
	}
	parse := func(s string) ([]int, error) {
		f := strings.Fields(s)
		if len(f) != n {
			return nil, fmt.Errorf("expected %d numbers", n)
		}
		res := make([]int, n)
		for i, x := range f {
			v, err := strconv.Atoi(x)
			if err != nil {
				return nil, fmt.Errorf("bad integer %q", x)
			}
			res[i] = v
		}
		return res, nil
	}
	v, err := parse(lines[0])
	if err != nil {
		return err
	}
	u, err := parse(lines[1])
	if err != nil {
		return err
	}
	return checkPaths(n, k, a, b, c, d, v, u)
}

func generateCase(rng *rand.Rand) (string, int, int, int, int, int, int) {
	n := rng.Intn(6) + 5     // 5..10
	k := rng.Intn(n) + n + 1 // >= n+1
	perm := rng.Perm(n)
	a := perm[0] + 1
	b := perm[1] + 1
	c := perm[2] + 1
	d := perm[3] + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, c, d))
	return sb.String(), n, k, a, b, c, d
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, k, a, b, c, d := generateCase(rng)
		if err := runCase(exe, in, n, k, a, b, c, d); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
