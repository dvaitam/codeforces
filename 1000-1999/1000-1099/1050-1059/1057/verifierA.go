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
	parents []int // 1-indexed, parents[1] unused
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

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 2 // 2..101
	p := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p[i] = rng.Intn(i-1) + 1
	}
	return testCase{parents: p}
}

func expected(tc testCase) string {
	n := len(tc.parents) - 1
	path := []int{}
	cur := n
	for cur != 1 {
		path = append(path, cur)
		cur = tc.parents[cur]
	}
	var sb strings.Builder
	sb.WriteString("1")
	for i := len(path) - 1; i >= 0; i-- {
		fmt.Fprintf(&sb, " %d", path[i])
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	n := len(tc.parents) - 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 2; i <= n; i++ {
		if i > 2 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", tc.parents[i])
	}
	b.WriteByte('\n')

	got, err := run(bin, b.String())
	if err != nil {
		return err
	}
	exp := expected(tc)
	if strings.TrimSpace(got) != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := make([]testCase, 0, 100)
	// deterministic small cases
	cases = append(cases, testCase{parents: []int{0, 0, 1}})    // n=2
	cases = append(cases, testCase{parents: []int{0, 0, 1, 1}}) // n=3 path 1->3
	cases = append(cases, testCase{parents: []int{0, 0, 1, 2}}) // n=3 path 1->2->3
	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
