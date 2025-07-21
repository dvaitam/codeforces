package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseB struct {
	input    string
	expected [][]int
}

func generateCase(rng *rand.Rand) testCaseB {
	n := rng.Intn(3) + 2 // 2..4
	sets := make([][]int, n)
	next := 1
	for i := 0; i < n; i++ {
		sz := rng.Intn(3) + 1 // 1..3
		sets[i] = make([]int, sz)
		for j := 0; j < sz; j++ {
			sets[i][j] = next
			next++
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			union := append(append([]int{}, sets[i]...), sets[j]...)
			sort.Ints(union)
			fmt.Fprintf(&b, "%d", len(union))
			for _, v := range union {
				fmt.Fprintf(&b, " %d", v)
			}
			fmt.Fprintln(&b)
		}
	}
	return testCaseB{input: b.String(), expected: sets}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	sets := make([][]int, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		k := 0
		fmt.Sscan(fields[0], &k)
		if len(fields)-1 != k {
			return nil, fmt.Errorf("invalid line %q", line)
		}
		arr := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Sscan(fields[i+1], &arr[i])
		}
		sort.Ints(arr)
		sets = append(sets, arr)
	}
	// sort sets lexicographically
	sort.Slice(sets, func(i, j int) bool {
		a, b := sets[i], sets[j]
		for x := 0; x < len(a) && x < len(b); x++ {
			if a[x] != b[x] {
				return a[x] < b[x]
			}
		}
		return len(a) < len(b)
	})
	return sets, nil
}

func equalSets(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("case %d invalid output: %v\n", i+1, err)
			os.Exit(1)
		}
		// sort expected sets as well
		exp := make([][]int, len(tc.expected))
		for j := range tc.expected {
			exp[j] = append([]int{}, tc.expected[j]...)
			sort.Ints(exp[j])
		}
		sort.Slice(exp, func(i, j int) bool {
			a, b := exp[i], exp[j]
			for x := 0; x < len(a) && x < len(b); x++ {
				if a[x] != b[x] {
					return a[x] < b[x]
				}
			}
			return len(a) < len(b)
		})
		if !equalSets(exp, got) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%v\ngot:%v\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
