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

type edge struct{ u, v int }

type testF struct {
	n     int
	bus   []edge
	train []edge
}

func genTreeEdges(rng *rand.Rand, offset, n int) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, edge{offset + u - 1, offset + i - 1})
	}
	return edges
}

func genTest(rng *rand.Rand) testF {
	n := rng.Intn(4) + 2 // 2..5
	bus := genTreeEdges(rng, 1, n)
	train := genTreeEdges(rng, n+1, n)
	return testF{n, bus, train}
}

func solveF(tc testF) (bool, []int) {
	n := tc.n
	trainSet := make(map[[2]int]bool)
	for _, e := range tc.train {
		a, b := e.u, e.v
		if a > b {
			a, b = b, a
		}
		trainSet[[2]int{a, b}] = true
	}
	perm := make([]int, n)
	used := make([]bool, n)
	nodes := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = n + 1 + i
	}
	var res []int
	var dfs func(int) bool
	dfs = func(pos int) bool {
		if pos == n {
			res = append([]int(nil), perm...)
			return true
		}
		for i, v := range nodes {
			if used[i] {
				continue
			}
			perm[pos] = v
			valid := true
			for _, e := range tc.bus {
				if e.u-1 == pos && used[e.v-1] {
					a := perm[pos]
					b := perm[e.v-1]
					if a > b {
						a, b = b, a
					}
					if trainSet[[2]int{a, b}] {
						valid = false
						break
					}
				} else if e.v-1 == pos && used[e.u-1] {
					a := perm[e.u-1]
					b := perm[pos]
					if a > b {
						a, b = b, a
					}
					if trainSet[[2]int{a, b}] {
						valid = false
						break
					}
				}
			}
			if valid {
				used[i] = true
				if dfs(pos + 1) {
					return true
				}
				used[i] = false
			}
		}
		return false
	}
	if dfs(0) {
		return true, res
	}
	return false, nil
}

func formatInput(tc testF) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for _, e := range tc.bus {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, e := range tc.train {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		input := formatInput(tc)
		ok, mapping := solveF(tc)
		var exp strings.Builder
		if !ok {
			exp.WriteString("No")
		} else {
			exp.WriteString("Yes\n")
			for j, v := range mapping {
				if j > 0 {
					exp.WriteByte(' ')
				}
				fmt.Fprintf(&exp, "%d", v)
			}
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp.String()) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp.String(), got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
