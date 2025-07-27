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
	pairs [][2]int
}

func buildCase(pairs [][2]int) testCase {
	return testCase{pairs: pairs}
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(10) + 1
	pairs := make([][2]int, t)
	for i := 0; i < t; i++ {
		a := rng.Intn(1_000_000_000)
		b := rng.Intn(1_000_000_000)
		pairs[i] = [2]int{a, b}
	}
	return buildCase(pairs)
}

func expected(tc testCase) []int {
	res := make([]int, len(tc.pairs))
	for i, p := range tc.pairs {
		res[i] = p[0] ^ p[1]
	}
	return res
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.pairs))
	for _, p := range tc.pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	exp := expected(tc)
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
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
	cases = append(cases, buildCase([][2]int{{6, 12}}))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			var sb strings.Builder
			fmt.Fprintf(&sb, "%d\n", len(tc.pairs))
			for _, p := range tc.pairs {
				fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
			}
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
