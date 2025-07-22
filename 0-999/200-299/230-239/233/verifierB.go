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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func sumDigits(x uint64) uint64 {
	var s uint64
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func integerSqrt(n uint64) uint64 {
	var low, high uint64 = 0, 2000000000
	for low < high {
		mid := (low + high + 1) >> 1
		if mid*mid <= n {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func expected(n uint64) int64 {
	var ans uint64
	found := false
	for s := uint64(1); s <= 200; s++ {
		D := s*s + 4*n
		r := integerSqrt(D)
		if r*r != D {
			continue
		}
		if r < s || (r-s)&1 != 0 {
			continue
		}
		x := (r - s) >> 1
		if x == 0 {
			continue
		}
		if sumDigits(x) != s {
			continue
		}
		if !found || x < ans {
			ans = x
			found = true
		}
	}
	if !found {
		return -1
	}
	return int64(ans)
}

func generateCases() []uint64 {
	cases := make([]uint64, 0, 100)
	for i := uint64(1); i <= 20; i++ {
		cases = append(cases, i)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		if rng.Intn(2) == 0 {
			x := uint64(rng.Int63n(1e9) + 1)
			n := x*x + sumDigits(x)*x
			cases = append(cases, n)
		} else {
			cases = append(cases, uint64(rng.Int63n(1e18)))
		}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ns := generateCases()
	for i, n := range ns {
		input := fmt.Sprintf("%d\n", n)
		outStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %s\n", i+1, outStr)
			os.Exit(1)
		}
		exp := expected(n)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
