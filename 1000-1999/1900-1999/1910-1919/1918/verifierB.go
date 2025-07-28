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

type testCase struct {
	a []int
	b []int
}

func solveCase(tc testCase) string {
	n := len(tc.a)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return tc.a[idx[i]] < tc.a[idx[j]] })
	var sb strings.Builder
	for i, id := range idx {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.a[id]))
	}
	sb.WriteByte('\n')
	for i, id := range idx {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.b[id]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	a := rng.Perm(n)
	b := rng.Perm(n)
	for i := 0; i < n; i++ {
		a[i]++
		b[i]++
	}
	return testCase{a: a, b: b}
}

func runCase(bin string, tc testCase) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solveCase(tc)
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := make([]testCase, 0, 100)
	cases = append(cases, testCase{a: []int{1, 2}, b: []int{2, 1}})
	for i := 0; i < 99; i++ {
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
