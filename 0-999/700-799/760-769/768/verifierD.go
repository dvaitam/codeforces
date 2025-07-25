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
	k       int
	q       int
	queries []int
	expect  []int
}

func expectedAnswers(k int, queries []int) []int {
	const maxP = 1000
	ans := make([]int, maxP+1)
	prev := make([]float64, k+1)
	curr := make([]float64, k+1)
	prev[0] = 1.0
	idx := 1
	n := 0
	kf := float64(k)
	for idx <= maxP {
		n++
		curr[0] = 0
		for j := 1; j <= k; j++ {
			curr[j] = prev[j]*float64(j)/kf + prev[j-1]*float64(k-j+1)/kf
		}
		prev, curr = curr, prev
		p := prev[k]
		for idx <= maxP && p >= float64(idx)/2000.0-1e-12 {
			ans[idx] = n
			idx++
		}
	}
	res := make([]int, len(queries))
	for i, pi := range queries {
		if pi < 1 {
			pi = 1
		} else if pi > 1000 {
			pi = 1000
		}
		res[i] = ans[pi]
	}
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	deterministic := []testCase{
		{k: 1, q: 3, queries: []int{1, 500, 1000}},
		{k: 3, q: 2, queries: []int{1000, 1}},
		{k: 5, q: 1, queries: []int{600}},
	}
	for _, tc := range deterministic {
		tc.expect = expectedAnswers(tc.k, tc.queries)
		tests = append(tests, tc)
	}
	for len(tests) < 100 {
		k := rng.Intn(10) + 1
		q := rng.Intn(5) + 1
		qs := make([]int, q)
		for i := 0; i < q; i++ {
			qs[i] = rng.Intn(1000) + 1
		}
		expect := expectedAnswers(k, qs)
		tests = append(tests, testCase{k: k, q: q, queries: qs, expect: expect})
	}
	return tests
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.k, tc.q))
		for _, v := range tc.queries {
			input.WriteString(fmt.Sprintf("%d\n", v))
		}
		expectParts := make([]string, len(tc.expect))
		for i2, v := range tc.expect {
			expectParts[i2] = strconv.Itoa(v)
		}
		expected := strings.Join(expectParts, "\n")
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n--- got:\n%s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
