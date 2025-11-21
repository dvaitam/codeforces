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
	n int
	a []int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
	}
	return testCase{n: n, a: a}
}

func bruteMinClones(a []int) int {
	n := len(a)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	best := n + 1
	var permute func(int)
	permute = func(pos int) {
		if pos == n {
			// check non-increasing weights
			for i := 1; i < n; i++ {
				if a[idx[i-1]] < a[idx[i]] {
					return
				}
			}
			active := make([]bool, n)
			parent := make([]int, n)
			for i := range parent {
				parent[i] = i
			}
			find := func(x int) int {
				for parent[x] != x {
					parent[x] = parent[parent[x]]
					x = parent[x]
				}
				return x
			}
			union := func(x, y int) {
				fx, fy := find(x), find(y)
				if fx != fy {
					parent[fy] = fx
				}
			}
			clones := 0
			for _, id := range idx {
				leftActive := id > 0 && active[id-1]
				rightActive := id+1 < n && active[id+1]
				if !leftActive && !rightActive {
					clones++
				} else if leftActive && rightActive {
					union(id-1, id)
					union(id, id+1)
				} else if leftActive {
					union(id-1, id)
				} else {
					union(id, id+1)
				}
				active[id] = true
			}
			if clones < best {
				best = clones
			}
			return
		}
		for i := pos; i < n; i++ {
			idx[pos], idx[i] = idx[i], idx[pos]
			permute(pos + 1)
			idx[pos], idx[i] = idx[i], idx[pos]
		}
	}
	permute(0)
	return best
}

func buildInput(cases []testCase) (string, []int) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	exp := make([]int, len(cases))
	for i, tc := range cases {
		fmt.Fprintln(&sb, tc.n)
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		exp[i] = bruteMinClones(tc.a)
	}
	return sb.String(), exp
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/2108C_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		got, err := strconv.Atoi(lines[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on case %d: %q\n", i+1, lines[i])
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\nn=%d a=%v\n", i+1, exp, got, cases[i].n, cases[i].a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
