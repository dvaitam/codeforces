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

type testCase struct {
	n        int
	a        []int
	b        []int
	bestLen  int
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

func bruteForceBestLen(a []int, b []int) int {
	gaps := []int{}
	for i, v := range a {
		if v == -1 {
			gaps = append(gaps, i)
		}
	}
	bestLen := -1

	var dfs func(idx int, used []bool)
	dfs = func(idx int, used []bool) {
		if idx == len(gaps) {
			cur := lisLength(a)
			if cur > bestLen {
				bestLen = cur
			}
			return
		}
		for i := 0; i < len(b); i++ {
			if used[i] {
				continue
			}
			used[i] = true
			a[gaps[idx]] = b[i]
			dfs(idx+1, used)
			used[i] = false
			a[gaps[idx]] = -1
		}
	}
	used := make([]bool, len(b))
	dfs(0, used)
	if bestLen == -1 {
		bestLen = lisLength(a)
	}
	return bestLen
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
	bl := bruteForceBestLen(append([]int(nil), a...), b)
	return testCase{n: n, a: a, b: b, bestLen: bl}
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

	// Parse output
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(parts))
	}
	result := make([]int, tc.n)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("bad output token %q", p)
		}
		result[i] = v
	}

	// Check that fixed positions are unchanged
	for i, v := range tc.a {
		if v != -1 && result[i] != v {
			return fmt.Errorf("position %d: fixed value %d changed to %d", i, v, result[i])
		}
	}

	// Check LIS length matches optimal
	gotLen := lisLength(result)
	if gotLen != tc.bestLen {
		return fmt.Errorf("expected LIS length %d got %d", tc.bestLen, gotLen)
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
