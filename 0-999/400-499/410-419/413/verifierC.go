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
	n, m     int
	prices   []int64
	auctions []int
}

func expected(tc testCase) int64 {
	isAuction := make([]bool, tc.n)
	for _, v := range tc.auctions {
		if v >= 1 && v <= tc.n {
			isAuction[v-1] = true
		}
	}
	var regSum int64
	var A []int64
	for i := 0; i < tc.n; i++ {
		if isAuction[i] {
			A = append(A, tc.prices[i])
		} else {
			regSum += tc.prices[i]
		}
	}
	sort.Slice(A, func(i, j int) bool { return A[i] < A[j] })
	totalSum := regSum
	for _, v := range A {
		totalSum += v
	}
	mA := len(A)
	PS := make([]int64, mA+1)
	for i := 1; i <= mA; i++ {
		PS[i] = PS[i-1] + A[i-1]
	}
	best := totalSum
	for d := 1; d <= mA; d++ {
		sumD := PS[d]
		S0 := totalSum - sumD
		if S0 <= 0 {
			continue
		}
		cur := S0
		ok := true
		for j := 1; j <= d; j++ {
			if cur <= A[j-1] {
				ok = false
				break
			}
			cur <<= 1
		}
		if !ok {
			continue
		}
		endScore := S0 << uint(d)
		if endScore > best {
			best = endScore
		}
	}
	return best
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(tc.prices[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(tc.auctions[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	m := rng.Intn(n + 1)
	prices := make([]int64, n)
	for i := 0; i < n; i++ {
		prices[i] = int64(rng.Intn(20) + 1)
	}
	auctions := make([]int, m)
	used := map[int]bool{}
	for i := 0; i < m; i++ {
		for {
			x := rng.Intn(n) + 1
			if !used[x] {
				auctions[i] = x
				used[x] = true
				break
			}
		}
	}
	return testCase{n: n, m: m, prices: prices, auctions: auctions}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, testCase{n: 1, m: 0, prices: []int64{5}, auctions: nil})
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
