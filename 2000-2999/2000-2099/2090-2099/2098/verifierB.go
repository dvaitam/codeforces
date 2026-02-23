package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func countHouses(arr []int, k int) int64 {
	n := len(arr)
	sort.Ints(arr)
	minM := n - k
	if minM%2 != 0 {
		minM++
	}
	if minM > n {
		return 1
	} else {
		c := minM / 2
		return int64(arr[n-c] - arr[c-1] + 1)
	}
}

type testCase struct {
	n  int
	k  int
	as []int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 1
	k := rng.Intn(n)
	as := make([]int, n)
	for i := 0; i < n; i++ {
		as[i] = rng.Intn(20) + 1
	}
	return testCase{n: n, k: k, as: as}
}

func buildInput(cases []testCase) (string, []int64) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	exp := make([]int64, len(cases))
	for i, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for j, v := range tc.as {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		exp[i] = countHouses(append([]int{}, tc.as...), tc.k)
	}
	return sb.String(), exp
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/2098B_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		gotStr := lines[i]
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on case %d: %q\n", i+1, gotStr)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\nn=%d k=%d arr=%v\n", i+1, exp, got, cases[i].n, cases[i].k, cases[i].as)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
