package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n, k int
	arr  []int
}

func solve(n, k int, arr []int) string {
	seen := make(map[int]bool)
	ans := make([]int, 0, k)
	for i, v := range arr {
		if !seen[v] {
			seen[v] = true
			if len(ans) < k {
				ans = append(ans, i+1)
			}
		}
	}
	var sb strings.Builder
	if len(ans) >= k {
		sb.WriteString("YES\n")
		for i := 0; i < k; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", ans[i])
		}
		sb.WriteByte('\n')
	} else {
		sb.WriteString("NO\n")
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", tc.n, tc.k)
	for i, v := range tc.arr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := strings.TrimSpace(solve(tc.n, tc.k, tc.arr))
	got := strings.TrimSpace(out.String())
	if expected != got {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var cases []testCase
	// some deterministic edge cases
	cases = append(cases, testCase{n: 1, k: 1, arr: []int{5}})
	cases = append(cases, testCase{n: 5, k: 2, arr: []int{1, 2, 1, 2, 3}})
	cases = append(cases, testCase{n: 5, k: 4, arr: []int{1, 1, 1, 1, 1}})

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		k := rng.Intn(n) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(100) + 1
		}
		cases = append(cases, testCase{n: n, k: k, arr: arr})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "input n=%d k=%d arr=%v\n", tc.n, tc.k, tc.arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
