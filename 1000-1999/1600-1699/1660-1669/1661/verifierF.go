package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
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

const inf int64 = 1<<63 - 1

func cost(L, p int64) int64 {
	if p <= 0 {
		return inf
	}
	q := L / p
	r := L % p
	qq := uint64(q * q)
	hi1, lo1 := bits.Mul64(qq, uint64(p-r))
	qq1 := uint64(q+1) * uint64(q+1)
	hi2, lo2 := bits.Mul64(qq1, uint64(r))
	lo, carry := bits.Add64(lo1, lo2, 0)
	hi, _ := bits.Add64(hi1, hi2, carry)
	if hi > 0 || lo > uint64(inf) {
		return inf
	}
	return int64(lo)
}

func delta(L, p int64) int64 {
	c1 := cost(L, p)
	c2 := cost(L, p+1)
	if c1 == inf || c2 == inf {
		return inf
	}
	return c1 - c2
}

func minParts(L, lam int64) int64 {
	low, high := int64(1), L
	for low < high {
		mid := (low + high) >> 1
		if delta(L, mid) <= lam {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func feasible(lengths []int64, lam, m int64) bool {
	var total int64
	for _, L := range lengths {
		p := minParts(L, lam)
		total += cost(L, p)
		if total > m {
			return false
		}
	}
	return total <= m
}

func solveCase(a []int64, m int64) int64 {
	n := len(a) - 1
	lengths := make([]int64, n)
	prev := int64(0)
	for i := 0; i < n; i++ {
		lengths[i] = a[i+1] - prev
		prev = a[i+1]
	}
	lo, hi := int64(0), int64(1e18)
	best := int64(0)
	for lo <= hi {
		mid := (lo + hi) >> 1
		if feasible(lengths, mid, m) {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	parts := make([]int64, len(lengths))
	var costNow int64
	for i, L := range lengths {
		p := minParts(L, best)
		parts[i] = p
		costNow += cost(L, p)
	}
	K := int64(0)
	for _, p := range parts {
		K += p - 1
	}
	leftover := m - costNow
	type item struct {
		d   int64
		idx int
	}
	pq := make([]item, 0)
	for i, L := range lengths {
		if parts[i] > 1 {
			d := delta(L, parts[i]-1)
			pq = append(pq, item{d: d, idx: i})
		}
	}
	sort.Slice(pq, func(i, j int) bool { return pq[i].d < pq[j].d })
	for len(pq) > 0 {
		it := pq[0]
		if it.d > leftover {
			break
		}
		leftover -= it.d
		K--
		idx := it.idx
		parts[idx]--
		if parts[idx] > 1 {
			d := delta(lengths[idx], parts[idx]-1)
			pq[0] = item{d: d, idx: idx}
		} else {
			pq = pq[1:]
		}
		sort.Slice(pq, func(i, j int) bool { return pq[i].d < pq[j].d })
	}
	return K
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// simple case n=1
	cases = append(cases, func() testCase {
		a := []int64{0, 1}
		m := int64(1)
		var sb strings.Builder
		sb.WriteString("1\n1\n1\n")
		exp := fmt.Sprintf("%d", solveCase(a, m))
		return testCase{input: sb.String(), expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(5) + 1
		a := make([]int64, n+1)
		cur := int64(0)
		for i := 1; i <= n; i++ {
			cur += int64(rng.Intn(100) + 1)
			a[i] = cur
		}
		m := cur + int64(rng.Intn(1000)+1)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", m))
		exp := fmt.Sprintf("%d", solveCase(a, m))
		cases = append(cases, testCase{input: sb.String(), expected: exp})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
