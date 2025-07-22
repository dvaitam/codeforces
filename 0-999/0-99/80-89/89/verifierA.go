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

type testCaseA struct {
	n int
	m int64
	k int64
	a []int64
}

func solveA(tc testCaseA) int64 {
	n := tc.n
	m := tc.m
	k := tc.k
	a := tc.a
	if n%2 == 0 {
		return 0
	}
	t := make([]int64, n)
	for i := 1; i < n; i++ {
		t[i] = a[i-1] + a[i] - t[i-1]
	}
	L := int64(-1 << 60)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			if v := -t[i]; v > L {
				L = v
			}
		}
	}
	x0 := a[0]
	maxRed := x0 - L
	if maxRed < 0 {
		maxRed = 0
	}
	cost := int64(n/2 + 1)
	per := m / cost
	if per <= 0 {
		return 0
	}
	total := per * k
	if total > maxRed {
		total = maxRed
	}
	return total
}

func runCaseA(bin string, tc testCaseA) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := solveA(tc)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func genCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(10) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(100000)
	}
	m := rng.Int63n(1000000000) + 1
	k := rng.Int63n(1000000000) + 1
	return testCaseA{n, m, k, a}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
