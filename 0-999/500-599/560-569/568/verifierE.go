package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n        int
	a        []int
	b        []int
	expected string
}

func lisLength(arr []int) int {
	d := make([]int, len(arr))
	length := 0
	for _, v := range arr {
		i := sort.Search(length, func(i int) bool { return d[i] >= v })
		d[i] = v
		if i == length {
			length++
		}
	}
	return length
}

func bruteForceFill(a []int, b []int) []int {
	// indexes of gaps
	gaps := []int{}
	for i, v := range a {
		if v == -1 {
			gaps = append(gaps, i)
		}
	}
	used := make([]bool, len(b))
	bestLen := -1
	var best []int

	var dfs func(idx int)
	dfs = func(idx int) {
		if idx == len(gaps) {
			tmp := make([]int, len(a))
			copy(tmp, a)
			cur := lisLength(tmp)
			if cur > bestLen {
				bestLen = cur
				best = append([]int(nil), tmp...)
			}
			return
		}
		for i := 0; i < len(b); i++ {
			if used[i] {
				continue
			}
			used[i] = true
			a[gaps[idx]] = b[i]
			dfs(idx + 1)
			used[i] = false
			a[gaps[idx]] = -1
		}
	}
	dfs(0)
	if best == nil {
		best = make([]int, len(a))
		copy(best, a)
		for i := range best {
			if best[i] == -1 {
				best[i] = b[0]
			}
		}
	}
	return best
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	a := make([]int, n)
	gaps := 0
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			a[i] = -1
			gaps++
		} else {
			a[i] = rng.Intn(10) + 1
		}
	}
	m := gaps + rng.Intn(3)
	if m == 0 {
		m = 1
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = rng.Intn(10) + 1
	}
	filled := bruteForceFill(append([]int(nil), a...), b)
	var sb strings.Builder
	for i, v := range filled {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return testCase{n: n, a: a, b: b, expected: sb.String()}
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.b)))
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.expected), strings.TrimSpace(out.String()))
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
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
