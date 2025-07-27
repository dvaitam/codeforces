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

type testCaseB struct {
	n int
	k int
	a []int64
}

func solveCaseB(n, k int, a []int64) int64 {
	b := append([]int64(nil), a...)
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	var sum int64
	for i := 0; i <= k; i++ {
		sum += b[i]
	}
	return sum
}

func buildInputB(tc testCaseB) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", tc.a[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCaseB(bin string, tc testCaseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputB(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solveCaseB(tc.n, tc.k, tc.a)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesB() []testCaseB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseB, 0, 100)
	// edge cases
	cases = append(cases, testCaseB{n: 2, k: 1, a: []int64{0, 0}})
	cases = append(cases, testCaseB{n: 3, k: 1, a: []int64{5, 3, 2}})
	for len(cases) < 100 {
		n := rng.Intn(10) + 2
		k := rng.Intn(n-1) + 1
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = rng.Int63n(1000)
		}
		cases = append(cases, testCaseB{n: n, k: k, a: arr})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesB()
	for i, tc := range cases {
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
