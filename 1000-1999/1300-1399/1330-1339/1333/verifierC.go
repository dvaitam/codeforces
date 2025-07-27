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
	input    string
	expected string
}

func solveCase(a []int64) string {
	prefix := make(map[int64]int)
	prefixSum := int64(0)
	prefix[0] = 0
	left := 0
	var res int64
	for i := 1; i <= len(a); i++ {
		prefixSum += a[i-1]
		if idx, ok := prefix[prefixSum]; ok && idx >= left {
			left = idx + 1
		}
		prefix[prefixSum] = i
		res += int64(i - left)
	}
	return fmt.Sprintf("%d\n", res)
}

func buildCase(a []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solveCase(a)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(11) - 5)
	}
	return buildCase(a)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
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
	cases := []testCase{
		buildCase([]int64{1, 2, -3}),
		buildCase([]int64{0, 0, 0}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
