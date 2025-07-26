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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isPossible(n, m int) bool {
	if m < n-1 {
		return false
	}
	if n < 1000 {
		cnt := 0
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if gcd(i, j) == 1 {
					cnt++
				}
			}
		}
		if cnt < m {
			return false
		}
	}
	return true
}

func genCaseD(rng *rand.Rand) (int, int) {
	n := rng.Intn(20) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	return n, m
}

func runCaseD(bin string, n, m int) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(scanner.Text())
	possible := isPossible(n, m)
	if !possible {
		if strings.ToLower(first) != "impossible" {
			return fmt.Errorf("expected Impossible got %s", first)
		}
		return nil
	}
	if strings.ToLower(first) != "possible" {
		return fmt.Errorf("expected Possible got %s", first)
	}
	// read edges
	type pair struct{ u, v int }
	edges := make(map[pair]struct{})
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		fa, fb := find(a), find(b)
		if fa != fb {
			parent[fa] = fb
		}
	}
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return fmt.Errorf("invalid edge line: %s", line)
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad ints")
		}
		if u < 1 || u > n || v < 1 || v > n || u == v {
			return fmt.Errorf("invalid vertices")
		}
		if gcd(u, v) != 1 {
			return fmt.Errorf("edge %d %d not coprime", u, v)
		}
		key := pair{u, v}
		if u > v {
			key = pair{v, u}
		}
		if _, ok := edges[key]; ok {
			return fmt.Errorf("duplicate edge")
		}
		edges[key] = struct{}{}
		union(u, v)
		count++
	}
	if count != m {
		return fmt.Errorf("expected %d edges got %d", m, count)
	}
	root := find(1)
	for i := 2; i <= n; i++ {
		if find(i) != root {
			return fmt.Errorf("graph not connected")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m := genCaseD(rng)
		if err := runCaseD(bin, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", i+1, err, n, m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
