package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	n, k     int
	arr      []int
	expected []string
}

func solveCase(tc testCaseE) []string {
	n, k := tc.n, tc.k
	arr := tc.arr
	res := make([]string, n-k+1)
	cnt := make(map[int]int)
	for i := 0; i < k; i++ {
		cnt[arr[i]]++
	}
	for i := 0; i <= n-k; i++ {
		max := math.MinInt64
		found := false
		for v, c := range cnt {
			if c == 1 {
				if !found || v > max {
					max = v
					found = true
				}
			}
		}
		if found {
			res[i] = fmt.Sprint(max)
		} else {
			res[i] = "Nothing"
		}
		if i == n-k {
			break
		}
		cnt[arr[i]]--
		if cnt[arr[i]] == 0 {
			delete(cnt, arr[i])
		}
		cnt[arr[i+k]]++
	}
	return res
}

func generateTests() []testCaseE {
	rng := rand.New(rand.NewSource(5))
	cases := make([]testCaseE, 100)
	for i := range cases {
		n := rng.Intn(30) + 1
		k := rng.Intn(n) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(11) - 5
		}
		tc := testCaseE{n: n, k: k, arr: arr}
		tc.expected = solveCase(tc)
		cases[i] = tc
	}
	return cases
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for _, v := range tc.arr {
			fmt.Fprintf(&sb, "%d\n", v)
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != len(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\n", i+1, len(tc.expected), len(lines))
			os.Exit(1)
		}
		for j, exp := range tc.expected {
			if strings.TrimSpace(lines[j]) != exp {
				fmt.Fprintf(os.Stderr, "case %d line %d expected %s got %s\n", i+1, j+1, exp, lines[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
