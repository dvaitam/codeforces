package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected []string
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func solveCase(names []string, edges [][2]int) []string {
	n := len(names)
	G := make([]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		G[a] |= 1 << b
		G[b] |= 1 << a
	}
	bestMask := 0
	best := 0
	total := 1 << n
	for mask := 0; mask < total; mask++ {
		cur := mask
		for j := 0; j < n; j++ {
			if mask&(1<<j) != 0 {
				cur &^= G[j]
			}
		}
		cnt := bits.OnesCount(uint(cur))
		if cnt > best {
			best = cnt
			bestMask = cur
		}
	}
	var res []string
	for i := 0; i < n; i++ {
		if bestMask&(1<<i) != 0 {
			res = append(res, names[i])
		}
	}
	if len(res) == 0 {
		res = []string{names[0]}
	} else {
		sort.Strings(res)
	}
	return res
}

func parseOutput(out string) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return nil, fmt.Errorf("read k: %v", err)
	}
	names := make([]string, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &names[i]); err != nil {
			return nil, fmt.Errorf("read name %d: %v", i+1, err)
		}
	}
	return names, nil
}

func runCase(bin string, tc testCase) error {
	out, err := run(bin, tc.input)
	if err != nil {
		return err
	}
	got, err := parseOutput(out)
	if err != nil {
		return err
	}
	if len(got) != len(tc.expected) {
		return fmt.Errorf("expected %d names got %d", len(tc.expected), len(got))
	}
	sort.Strings(got)
	for i := range got {
		if got[i] != tc.expected[i] {
			return fmt.Errorf("expected %v got %v", tc.expected, got)
		}
	}
	return nil
}

func buildCase(names []string, edges [][2]int) testCase {
	var sb strings.Builder
	n := len(names)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, name := range names {
		sb.WriteString(name + "\n")
	}
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%s %s\n", names[e[0]], names[e[1]]))
	}
	return testCase{input: sb.String(), expected: solveCase(names, edges)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("p%d", i+1)
	}
	var edges [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(3) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return buildCase(names, edges)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase([]string{"a"}, nil))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
