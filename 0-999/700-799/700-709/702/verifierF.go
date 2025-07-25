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

type Shirt struct {
	c, q int64
}

type testCase struct {
	input    string
	expected string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func solve(shirts []Shirt, budgets []int64) []int {
	sort.Slice(shirts, func(i, j int) bool {
		if shirts[i].q != shirts[j].q {
			return shirts[i].q > shirts[j].q
		}
		return shirts[i].c < shirts[j].c
	})
	prefix := make([]int64, len(shirts)+1)
	for i := 0; i < len(shirts); i++ {
		prefix[i+1] = prefix[i] + shirts[i].c
	}
	res := make([]int, len(budgets))
	for i, b := range budgets {
		lo, hi := 0, len(shirts)
		for lo < hi {
			mid := (lo + hi + 1) >> 1
			if prefix[mid] <= b {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		res[i] = lo
	}
	return res
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	shirts := make([]Shirt, n)
	for i := range shirts {
		shirts[i] = Shirt{c: int64(rng.Intn(10) + 1), q: int64(rng.Intn(10) + 1)}
	}
	k := rng.Intn(6) + 1
	budgets := make([]int64, k)
	for i := range budgets {
		budgets[i] = int64(rng.Intn(50) + 1)
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, s := range shirts {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d %d", s.c, s.q))
		if i != n-1 {
			input.WriteByte(' ')
		}
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", k))
	for i, b := range budgets {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(b, 10))
	}
	input.WriteByte('\n')
	answers := solve(shirts, budgets)
	var expected strings.Builder
	for i, v := range answers {
		if i > 0 {
			expected.WriteByte(' ')
		}
		expected.WriteString(strconv.Itoa(v))
	}
	expected.WriteByte('\n')
	return testCase{input: input.String(), expected: strings.TrimSpace(expected.String())}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, genRandomCase(rand.New(rand.NewSource(42))))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
