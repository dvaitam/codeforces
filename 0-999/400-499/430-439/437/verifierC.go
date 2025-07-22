package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(n, m int, vals []int, edges [][2]int) int {
	ans := 0
	for _, e := range edges {
		x := e[0]
		y := e[1]
		if vals[x-1] < vals[y-1] {
			ans += vals[x-1]
		} else {
			ans += vals[y-1]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(100)
	}
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, [2]int{a, b})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	exp := fmt.Sprintf("%d", solveC(n, m, vals, edges))
	return input, exp
}

func manualCase() (string, string) {
	n := 4
	m := 3
	vals := []int{10, 20, 30, 40}
	edges := [][2]int{{1, 2}, {2, 3}, {4, 1}}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "10 20 30 40\n")
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	exp := fmt.Sprintf("%d", solveC(n, m, vals, edges))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	manualIn, manualExp := manualCase()
	cases = append(cases, [2]string{manualIn, manualExp})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for idx, tc := range cases {
		out, err := runBinary(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
