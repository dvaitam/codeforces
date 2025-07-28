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

type testCaseA struct {
	n int64
	m int64
}

func expectedA(n, m int64) int64 {
	if n == 1 && m == 1 {
		return 0
	}
	if n-1 < m-1 {
		return n + m - 1 + (n - 1)
	}
	return n + m - 1 + (m - 1)
}

func genCaseA(rng *rand.Rand) []testCaseA {
	t := rng.Intn(5) + 1
	cases := make([]testCaseA, t)
	for i := 0; i < t; i++ {
		cases[i].n = rng.Int63n(1e5) + 1
		cases[i].m = rng.Int63n(1e5) + 1
	}
	return cases
}

func runCaseA(bin string, tc []testCaseA) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tc)))
	for _, c := range tc {
		input.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(tc) {
		return fmt.Errorf("expected %d numbers got %d", len(tc), len(fields))
	}
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if val != expectedA(tc[i].n, tc[i].m) {
			return fmt.Errorf("case %d expected %d got %d", i+1, expectedA(tc[i].n, tc[i].m), val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			var input strings.Builder
			input.WriteString(fmt.Sprintf("%d\n", len(tc)))
			for _, c := range tc {
				input.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
			}
			fmt.Fprint(os.Stderr, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
