package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	n int
	k int64
	m int64
	a []int64
}

func solve(n int, k, m int64, a []int64) string {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var sum int64
	for _, v := range a {
		sum += v
	}
	ans := float64(sum) / float64(n)
	var suffix int64
	for i := n - 1; i >= 0; i-- {
		suffix += a[i]
		s := int64(n - i)
		ops := m - int64(i)
		if ops < 0 {
			continue
		}
		add := k * s
		if add > ops {
			add = ops
		}
		cand := float64(suffix+add) / float64(s)
		if cand > ans {
			ans = cand
		}
	}
	return fmt.Sprintf("%.10f", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		k := int64(rng.Intn(5) + 1)
		m := int64(rng.Intn(10) + 1)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = int64(rng.Intn(10) + 1)
		}
		tests = append(tests, test{n, k, m, a})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.k, tc.m)
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := solve(tc.n, tc.k, tc.m, append([]int64(nil), tc.a...))
		got, err := run(binary, input)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
