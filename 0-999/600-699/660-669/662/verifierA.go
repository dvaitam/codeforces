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
	a []uint64
	b []uint64
}

func solveA(tc testCaseA) string {
	base := uint64(0)
	diffs := make([]uint64, tc.n)
	for i := 0; i < tc.n; i++ {
		base ^= tc.a[i]
		diffs[i] = tc.a[i] ^ tc.b[i]
	}
	var basis [64]uint64
	rank := 0
	for _, d := range diffs {
		x := d
		for j := 63; j >= 0; j-- {
			if (x>>uint(j))&1 == 0 {
				continue
			}
			if basis[j] == 0 {
				basis[j] = x
				rank++
				break
			}
			x ^= basis[j]
		}
	}
	inSpan := func(v uint64) bool {
		x := v
		for j := 63; j >= 0; j-- {
			if (x>>uint(j))&1 == 0 {
				continue
			}
			if basis[j] == 0 {
				return false
			}
			x ^= basis[j]
		}
		return true
	}(base)
	if !inSpan {
		return "1/1"
	}
	denom := uint64(1) << uint(rank)
	numer := denom - 1
	return fmt.Sprintf("%d/%d", numer, denom)
}

func runCaseA(bin string, tc testCaseA) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveA(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func genCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(10) + 1
	a := make([]uint64, n)
	b := make([]uint64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Uint64() % 1000
		b[i] = rng.Uint64() % 1000
	}
	return testCaseA{n, a, b}
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
