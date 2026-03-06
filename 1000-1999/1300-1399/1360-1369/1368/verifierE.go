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

type testCase struct {
	n, m  int
	edges [][2]int
}

func genTests(rng *rand.Rand) (string, []testCase) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	cases := make([]testCase, t)
	for ci := 0; ci < t; ci++ {
		n := rng.Intn(8) + 2
		outDeg := make([]int, n+1)
		var edges [][2]int
		maxAttempts := n * n
		for a := 0; a < maxAttempts && len(edges) < 2*n; a++ {
			x := rng.Intn(n-1) + 1
			y := rng.Intn(n-x) + x + 1
			if outDeg[x] >= 2 {
				continue
			}
			dup := false
			for _, e := range edges {
				if e[0] == x && e[1] == y {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			outDeg[x]++
			edges = append(edges, [2]int{x, y})
		}
		m := len(edges)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		cases[ci] = testCase{n, m, edges}
	}
	return sb.String(), cases
}

func verify(tc testCase, scanner *bufio.Scanner) error {
	if !scanner.Scan() {
		return fmt.Errorf("missing k line")
	}
	k, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	if 7*k > 4*tc.n {
		return fmt.Errorf("too many closed: %d (need 7*k <= 4*%d = %d)", k, tc.n, 4*tc.n)
	}

	closed := make(map[int]bool)
	if k > 0 {
		if !scanner.Scan() {
			return fmt.Errorf("missing closed spots line")
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) != k {
			return fmt.Errorf("expected %d spots, got %d tokens", k, len(parts))
		}
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				return fmt.Errorf("bad spot: %v", err)
			}
			if v < 1 || v > tc.n {
				return fmt.Errorf("spot %d out of range [1,%d]", v, tc.n)
			}
			if closed[v] {
				return fmt.Errorf("duplicate spot %d", v)
			}
			closed[v] = true
		}
	} else {
		scanner.Scan()
	}

	hasIn := make([]bool, tc.n+1)
	hasOut := make([]bool, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		if closed[u] || closed[v] {
			continue
		}
		hasOut[u] = true
		hasIn[v] = true
	}
	for v := 1; v <= tc.n; v++ {
		if !closed[v] && hasIn[v] && hasOut[v] {
			return fmt.Errorf("vertex %d has both in and out edges — path of length >=2 exists", v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input, cases := genTests(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\ninput:\n%s", i+1, err, stderr.String(), input)
			os.Exit(1)
		}
		sc := bufio.NewScanner(strings.NewReader(stdout.String()))
		for ci, tc := range cases {
			if err := verify(tc, sc); err != nil {
				fmt.Fprintf(os.Stderr, "case %d subcase %d: %v\ninput:\n%sgot:\n%s\n", i+1, ci+1, err, input, stdout.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
