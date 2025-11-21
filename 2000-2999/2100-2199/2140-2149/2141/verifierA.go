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

type sofa struct {
	price int
	idx   int
}

type testCase struct {
	n   int
	arr []int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	// generate distinct prices between 1 and 100
	perm := rng.Perm(100)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = perm[i] + 1
	}
	return testCase{n: n, arr: arr}
}

func expectedBad(arr []int) []int {
	n := len(arr)
	s := make([]sofa, n)
	for i, v := range arr {
		s[i] = sofa{price: v, idx: i + 1}
	}
	sort.Slice(s, func(i, j int) bool { return s[i].price < s[j].price })
	sum := 0
	var res []int
	for i := 1; i < n; i++ {
		sum += s[i-1].price
		if sum < s[i].price {
			res = append(res, s[i].idx)
		}
	}
	sort.Ints(res)
	return res
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintln(&sb, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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

func parseAndValidate(tokens []string, idx *int, expected []int) error {
	if *idx >= len(tokens) {
		return fmt.Errorf("not enough output tokens")
	}
	k, err := strconv.Atoi(tokens[*idx])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid k: %q", tokens[*idx])
	}
	*idx++
	if k != len(expected) {
		return fmt.Errorf("k mismatch: got %d expected %d", k, len(expected))
	}
	got := make([]int, k)
	for i := 0; i < k; i++ {
		if *idx >= len(tokens) {
			return fmt.Errorf("not enough indices")
		}
		val, err := strconv.Atoi(tokens[*idx])
		if err != nil {
			return fmt.Errorf("invalid index: %q", tokens[*idx])
		}
		got[i] = val
		*idx++
	}
	if !sort.IntsAreSorted(got) {
		return fmt.Errorf("indices not sorted ascending: %v", got)
	}
	if len(got) != len(expected) {
		return fmt.Errorf("length mismatch after parsing")
	}
	for i := range got {
		if got[i] != expected[i] {
			return fmt.Errorf("indices mismatch: got %v expected %v", got, expected)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/2141A_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input := buildInput(cases)
	expectedAll := make([][]int, len(cases))
	for i, tc := range cases {
		expectedAll[i] = expectedBad(tc.arr)
	}
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(output)
	ptr := 0
	for i, exp := range expectedAll {
		if err := parseAndValidate(tokens, &ptr, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nn=%d arr=%v\n", i+1, err, cases[i].n, cases[i].arr)
			os.Exit(1)
		}
	}
	if ptr != len(tokens) {
		fmt.Fprintf(os.Stderr, "extra output tokens detected\n")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
