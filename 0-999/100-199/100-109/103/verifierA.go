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
	input  string
	expect int64
}

func solveCase(a []int64) int64 {
	ans := int64(len(a))
	for i, v := range a {
		ans += (v - 1) * int64(i+1)
	}
	return ans
}

func buildCase(a []int64) testCase {
	var sb strings.Builder
	n := len(a)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", a[i])
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expect: solveCase(a)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(100) + 1
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
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, got)
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

	var cases []testCase
	// fixed edge cases
	cases = append(cases, buildCase([]int64{1}))
	cases = append(cases, buildCase([]int64{5, 3, 2}))
	cases = append(cases, buildCase([]int64{1, 1, 1, 1}))

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
