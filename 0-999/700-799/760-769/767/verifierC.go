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

type testCase struct {
	n      int
	parent []int
	val    []int64
}

func solve(n int, parent []int, val []int64) string {
	children := make([][]int, n+1)
	var root int
	var total int64
	for i := 1; i <= n; i++ {
		if parent[i] == 0 {
			root = i
		} else {
			children[parent[i]] = append(children[parent[i]], i)
		}
		total += val[i]
	}
	if total%3 != 0 {
		return "-1"
	}
	target := total / 3
	sum := make([]int64, n+1)
	ans1, ans2 := 0, 0
	type frame struct {
		node int
		done bool
	}
	stack := []frame{{node: root, done: false}}
	for len(stack) > 0 && (ans1 == 0 || ans2 == 0) {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if f.done {
			sum[f.node] = val[f.node]
			for _, c := range children[f.node] {
				if sum[c] == target {
					if ans1 == 0 {
						ans1 = c
					} else if ans2 == 0 {
						ans2 = c
					}
				} else {
					sum[f.node] += sum[c]
				}
			}
		} else {
			stack = append(stack, frame{node: f.node, done: true})
			for _, c := range children[f.node] {
				stack = append(stack, frame{node: c, done: false})
			}
		}
	}
	if ans1 == 0 || ans2 == 0 {
		return "-1"
	}
	return fmt.Sprintf("%d %d", ans1, ans2)
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 1; i <= tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.parent[i], tc.val[i]))
	}
	return sb.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 2
	parent := make([]int, n+1)
	val := make([]int64, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	for i := 1; i <= n; i++ {
		val[i] = int64(rng.Intn(11) - 5)
	}
	return testCase{n: n, parent: parent, val: val}
}

func runProgram(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.n, append([]int(nil), tc.parent...), append([]int64(nil), tc.val...))
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{randomCase(rng)}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
