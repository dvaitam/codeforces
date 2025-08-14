package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	absBin, err := filepath.Abs(bin)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(absBin)
	inFile := filepath.Join(dir, "input.txt")
	outFile := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(inFile, []byte(input), 0o644); err != nil {
		return "", err
	}
	defer os.Remove(inFile)
	defer os.Remove(outFile)

	var cmd *exec.Cmd
	if strings.HasSuffix(absBin, ".go") {
		cmd = exec.Command("go", "run", absBin)
	} else {
		cmd = exec.Command(absBin)
	}
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	res := strings.TrimSpace(out.String())
	if res == "" {
		if data, err := os.ReadFile(outFile); err == nil {
			res = strings.TrimSpace(string(data))
		}
	}
	return res, nil
}

func bfsDistances(n, m int, starts [][2]int) ([]int, int) {
	total := n * m
	dist := make([]int, total)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, total)
	for _, st := range starts {
		x, y := st[0]-1, st[1]-1
		idx := x*m + y
		if dist[idx] == -1 {
			dist[idx] = 0
			q = append(q, idx)
		}
	}
	for head := 0; head < len(q); head++ {
		idx := q[head]
		r := idx / m
		c := idx % m
		d := dist[idx] + 1
		if r > 0 {
			ni := (r-1)*m + c
			if dist[ni] == -1 {
				dist[ni] = d
				q = append(q, ni)
			}
		}
		if r+1 < n {
			ni := (r+1)*m + c
			if dist[ni] == -1 {
				dist[ni] = d
				q = append(q, ni)
			}
		}
		if c > 0 {
			ni := r*m + c - 1
			if dist[ni] == -1 {
				dist[ni] = d
				q = append(q, ni)
			}
		}
		if c+1 < m {
			ni := r*m + c + 1
			if dist[ni] == -1 {
				dist[ni] = d
				q = append(q, ni)
			}
		}
	}
	maxDist := 0
	for _, v := range dist {
		if v > maxDist {
			maxDist = v
		}
	}
	return dist, maxDist
}

type testCase struct {
	input   string
	n, m    int
	dist    []int
	maxDist int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	k := rng.Intn(min(3, n*m)) + 1
	starts := make([][2]int, k)
	used := make(map[[2]int]bool)
	for i := 0; i < k; i++ {
		for {
			x := rng.Intn(n) + 1
			y := rng.Intn(m) + 1
			if !used[[2]int{x, y}] {
				used[[2]int{x, y}] = true
				starts[i] = [2]int{x, y}
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d\n", k)
	for i, st := range starts {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d %d", st[0], st[1])
	}
	sb.WriteByte('\n')
	dist, maxDist := bfsDistances(n, m, starts)
	return testCase{input: sb.String(), n: n, m: m, dist: dist, maxDist: maxDist}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// simple deterministic case
	cases = append(cases, testCase{
		input:   "1 1\n1\n1 1\n",
		n:       1,
		m:       1,
		dist:    []int{0},
		maxDist: 0,
	})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var x, y int
		if _, err := fmt.Sscan(out, &x, &y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output format: %v\noutput: %q\ninput:\n%s", i+1, err, out, tc.input)
			os.Exit(1)
		}
		if x < 1 || x > tc.n || y < 1 || y > tc.m {
			fmt.Fprintf(os.Stderr, "case %d failed: coordinates out of range: %d %d\ninput:\n%s", i+1, x, y, tc.input)
			os.Exit(1)
		}
		idx := (x-1)*tc.m + (y - 1)
		if tc.dist[idx] != tc.maxDist {
			// find one valid answer for message
			expIdx := 0
			for j, d := range tc.dist {
				if d == tc.maxDist {
					expIdx = j
					break
				}
			}
			ex := expIdx/tc.m + 1
			ey := expIdx%tc.m + 1
			fmt.Fprintf(os.Stderr, "case %d failed: expected distance %d (e.g., %d %d) got %d %d\ninput:\n%s", i+1, tc.maxDist, ex, ey, x, y, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
