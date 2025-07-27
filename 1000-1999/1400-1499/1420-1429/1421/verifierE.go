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
	a []int64
}

func buildCase(a []int64) testCase { return testCase{a: a} }

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(2_000_000_001) - 1_000_000_000
	}
	return buildCase(arr)
}

func solve(tc testCase) int64 {
	n := len(tc.a)
	if n == 1 {
		return tc.a[0]
	}
	if n == 2 {
		return -(tc.a[0] + tc.a[1])
	}
	var sumAbs int64
	minAbs := tc.a[0]
	if minAbs < 0 {
		minAbs = -minAbs
	}
	for _, v := range tc.a {
		val := v
		if val < 0 {
			val = -val
		}
		sumAbs += val
		if val < minAbs {
			minAbs = val
		}
	}
	if n%2 == 1 {
		return sumAbs - 2*minAbs
	}
	return sumAbs
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.a))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteString("\n")
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solve(tc)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int64{5, 6, 7, 8}))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			var sb strings.Builder
			fmt.Fprintf(&sb, "%d\n", len(tc.a))
			for i, v := range tc.a {
				if i > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", v)
			}
			sb.WriteString("\n")
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
