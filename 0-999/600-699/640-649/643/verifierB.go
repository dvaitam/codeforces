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

type testCase struct {
	n, k       int
	a, b, c, d int
	input      string
	possible   bool
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 2                   // 2..9 (some invalid when <=4)
	k := rng.Intn(2*n-2-(n-1)+1) + (n - 1) // range [n-1, 2n-2]
	a := rng.Intn(n) + 1
	b := rng.Intn(n-1) + 1
	if b >= a {
		b++
	}
	c := rng.Intn(n-2) + 1
	if c >= a {
		c++
	}
	if c >= b {
		c++
	}
	d := rng.Intn(n-3) + 1
	if d >= a {
		d++
	}
	if d >= b {
		d++
	}
	if d >= c {
		d++
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, c, d)

	possible := !(k <= n || n <= 4)

	return testCase{n, k, a, b, c, d, sb.String(), possible}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseLine(line string, n int) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int, n)
	used := make([]bool, n+1)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		if v < 1 || v > n {
			return nil, fmt.Errorf("value out of range")
		}
		if used[v] {
			return nil, fmt.Errorf("number %d repeated", v)
		}
		used[v] = true
		res[i] = v
	}
	return res, nil
}

func validate(tc testCase, out string) error {
	out = strings.TrimSpace(out)
	if !tc.possible {
		if out != "-1" {
			return fmt.Errorf("expected -1 for impossible case")
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected two lines in output")
	}
	path1, err := parseLine(lines[0], tc.n)
	if err != nil {
		return fmt.Errorf("first line: %v", err)
	}
	path2, err := parseLine(lines[1], tc.n)
	if err != nil {
		return fmt.Errorf("second line: %v", err)
	}
	if path1[0] != tc.a || path1[tc.n-1] != tc.b {
		return fmt.Errorf("first line should start with %d and end with %d", tc.a, tc.b)
	}
	if path2[0] != tc.c || path2[tc.n-1] != tc.d {
		return fmt.Errorf("second line should start with %d and end with %d", tc.c, tc.d)
	}
	// edges
	edges := make(map[[2]int]struct{})
	addEdge := func(x, y int) {
		if x > y {
			x, y = y, x
		}
		edges[[2]int{x, y}] = struct{}{}
	}
	for i := 0; i < tc.n-1; i++ {
		addEdge(path1[i], path1[i+1])
	}
	for i := 0; i < tc.n-1; i++ {
		addEdge(path2[i], path2[i+1])
	}
	if _, ok := edges[[2]int{min(tc.a, tc.b), max(tc.a, tc.b)}]; ok {
		return fmt.Errorf("edge a-b should not exist")
	}
	if _, ok := edges[[2]int{min(tc.c, tc.d), max(tc.c, tc.d)}]; ok {
		return fmt.Errorf("edge c-d should not exist")
	}
	if len(edges) > tc.k {
		return fmt.Errorf("too many edges: got %d expect <= %d", len(edges), tc.k)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := validate(c, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
