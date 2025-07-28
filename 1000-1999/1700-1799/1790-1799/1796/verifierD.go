package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	n int
	k int
	x int64
	a []int64
}

func maxSubarray(arr []int64) int64 {
	maxSum, cur := int64(0), int64(0)
	for _, v := range arr {
		cur += v
		if cur > maxSum {
			maxSum = cur
		}
		if cur < 0 {
			cur = 0
		}
	}
	return maxSum
}

func solve(tc TestCase) int64 {
	n := tc.n
	k := tc.k
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	best := int64(0)
	var dfs func(pos, chosen int, mask int)
	dfs = func(pos, chosen int, mask int) {
		if pos == n {
			if chosen != k {
				return
			}
			b := make([]int64, n)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					b[i] = tc.a[i] + tc.x
				} else {
					b[i] = tc.a[i] - tc.x
				}
			}
			val := maxSubarray(b)
			if val > best {
				best = val
			}
			return
		}
		if chosen < k {
			dfs(pos+1, chosen+1, mask|(1<<pos))
		}
		if n-pos-1 >= k-chosen {
			dfs(pos+1, chosen, mask)
		}
	}
	dfs(0, 0, 0)
	return best
}

func genTests() []TestCase {
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := 4 + i%3
		k := i % (n + 1)
		x := int64(i%5 - 2)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = int64((i+j)%5 - 2)
		}
		tests = append(tests, TestCase{n, k, x, a})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.k, tc.x)
		for i, v := range tc.a {
			if i+1 == len(tc.a) {
				fmt.Fprintf(&input, "%d\n", v)
			} else {
				fmt.Fprintf(&input, "%d ", v)
			}
		}
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = fmt.Sprintf("%d", solve(tc))
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expected) {
		fmt.Printf("wrong number of lines: got %d want %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, got := range lines {
		got = strings.TrimSpace(got)
		if got != expected[i] {
			fmt.Printf("test %d failed expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
