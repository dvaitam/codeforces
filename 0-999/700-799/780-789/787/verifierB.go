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

type testB struct {
	n, m   int
	groups [][]int
}

func solveB(n, m int, groups [][]int) string {
	for _, g := range groups {
		seen := make(map[int]bool, len(g))
		ok := true
		for _, v := range g {
			if seen[-v] {
				ok = false
			}
			seen[v] = true
		}
		if ok {
			return "YES\n"
		}
	}
	return "NO\n"
}

func runCase(bin string, input string, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func (tc testB) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, g := range tc.groups {
		sb.WriteString(fmt.Sprintf("%d", len(g)))
		for _, v := range g {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testB{
		{1, 1, [][]int{{1}}},
		{2, 1, [][]int{{-1, 2}}},
		{2, 2, [][]int{{1, -1}, {2}}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		groups := make([][]int, m)
		for j := 0; j < m; j++ {
			k := rng.Intn(5) + 1
			groups[j] = make([]int, k)
			for t := 0; t < k; t++ {
				v := rng.Intn(n) + 1
				if rng.Intn(2) == 0 {
					v = -v
				}
				groups[j][t] = v
			}
		}
		cases = append(cases, testB{n, m, groups})
	}
	for i, tc := range cases {
		input := tc.Input()
		expected := solveB(tc.n, tc.m, tc.groups)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
