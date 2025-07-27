package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// BIT structure for prefix sums

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
	for x := i + 1; x <= b.n; x += x & -x {
		b.tree[x] += v
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for x := i + 1; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}

func solveCaseE(s string) int64 {
	n := len(s)
	pos := make([][]int, 26)
	for i := 0; i < n; i++ {
		c := s[i] - 'a'
		pos[c] = append(pos[c], i)
	}
	P := make([]int, n)
	ptr := make([]int, 26)
	idx := 0
	for i := n - 1; i >= 0; i-- {
		c := s[i] - 'a'
		P[idx] = pos[c][ptr[c]]
		ptr[c]++
		idx++
	}
	bit := NewBIT(n)
	var inv int64
	for i, v := range P {
		cnt := bit.Sum(v)
		inv += int64(i - cnt)
		bit.Add(v, 1)
	}
	return inv
}

type testCaseE struct {
	s string
}

func buildInputE(s string) string {
	return fmt.Sprintf("%d\n%s\n", len(s), s)
}

func runCaseE(bin string, tc testCaseE) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputE(tc.s))
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
	expect := solveCaseE(tc.s)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesE() []testCaseE {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseE, 0, 100)
	cases = append(cases, testCaseE{s: "a"}, testCaseE{s: "ab"}, testCaseE{s: "ba"})
	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(5)))
		}
		cases = append(cases, testCaseE{s: sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesE()
	for i, tc := range cases {
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (s=%s)\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
