package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseB struct {
	n int
	k int64
	b int64
	s int64
}

func genCaseB(rng *rand.Rand) []testCaseB {
	t := rng.Intn(5) + 1
	cases := make([]testCaseB, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		k := int64(rng.Intn(10) + 1)
		b := int64(rng.Intn(10))
		base := b * k
		maxS := base + int64(n)*(k-1)
		var s int64
		if rng.Intn(2) == 0 {
			s = base + int64(rng.Intn(int(maxS-base+1)))
		} else {
			s = maxS + int64(rng.Intn(10)+1) // invalid
		}
		cases[i] = testCaseB{n: n, k: k, b: b, s: s}
	}
	return cases
}

func runCaseB(bin string, tcs []testCaseB) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tcs)))
	for _, tc := range tcs {
		input.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.k, tc.b, tc.s))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(tcs) {
		return fmt.Errorf("expected %d lines got %d", len(tcs), len(lines))
	}
	for i, line := range lines {
		tc := tcs[i]
		base := tc.b * tc.k
		maxSum := base + int64(tc.n)*(tc.k-1)
		possible := tc.s >= base && tc.s <= maxSum

		fields := strings.Fields(line)
		if !possible {
			if strings.TrimSpace(line) != "-1" {
				return fmt.Errorf("case %d expected -1 got %s", i+1, strings.TrimSpace(line))
			}
			continue
		}
		if strings.TrimSpace(line) == "-1" {
			return fmt.Errorf("case %d expected solution, got -1", i+1)
		}

		if len(fields) != tc.n {
			return fmt.Errorf("case %d expected %d numbers got %d", i+1, tc.n, len(fields))
		}

		var sum int64
		var beauty int64
		for _, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return fmt.Errorf("case %d invalid int %q", i+1, f)
			}
			if val < 0 {
				return fmt.Errorf("case %d negative value %d", i+1, val)
			}
			sum += val
			beauty += val / tc.k
		}
		if sum != tc.s {
			return fmt.Errorf("case %d sum mismatch expected %d got %d", i+1, tc.s, sum)
		}
		if beauty != tc.b {
			return fmt.Errorf("case %d beauty mismatch expected %d got %d", i+1, tc.b, beauty)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genCaseB(rng)
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			var inp strings.Builder
			inp.WriteString(fmt.Sprintf("%d\n", len(tc)))
			for _, c := range tc {
				inp.WriteString(fmt.Sprintf("%d %d %d %d\n", c.n, c.k, c.b, c.s))
			}
			fmt.Fprint(os.Stderr, inp.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
