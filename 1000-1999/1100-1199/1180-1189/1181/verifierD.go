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

func solve(n, m int, a []int, qs []uint64) []int {
	maxK := uint64(0)
	for _, k := range qs {
		if k > maxK {
			maxK = k
		}
	}
	schedule := make([]int, maxK+1)
	counts := make([]int, m+1)
	for i := 0; i < n; i++ {
		schedule[i+1] = a[i]
		counts[a[i]]++
	}
	for year := n + 1; uint64(year) <= maxK; year++ {
		best := 1
		for c := 1; c <= m; c++ {
			if counts[c] < counts[best] || (counts[c] == counts[best] && c < best) {
				best = c
			}
		}
		schedule[year] = best
		counts[best]++
	}
	ans := make([]int, len(qs))
	for i, k := range qs {
		ans[i] = schedule[k]
	}
	return ans
}

func runCase(bin string, n, m, q int, a []int, ks []uint64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", ks[i]))
	}
	input := sb.String()
	expectAns := solve(n, m, a, ks)
	var expect strings.Builder
	for i, v := range expectAns {
		if i > 0 {
			expect.WriteByte('\n')
		}
		expect.WriteString(fmt.Sprintf("%d", v))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect.String() {
		return fmt.Errorf("expected %q got %q", expect.String(), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		m := rng.Intn(4) + 1
		q := rng.Intn(3) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(m) + 1
		}
		ks := make([]uint64, q)
		for j := 0; j < q; j++ {
			ks[j] = uint64(rng.Intn(30) + n + 1)
		}
		if err := runCase(bin, n, m, q, a, ks); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
