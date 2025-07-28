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

type testCase struct {
	n, m int64
}

func solve(n, m int64) string {
	return fmt.Sprintf("%d", n*(m/2))
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Int63n(10000-2+1) + 2
	m := rng.Int63n(10000-2+1) + 2
	return testCase{n, m}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solve(tc.n, tc.m)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]

	// deterministic case
	cases := []testCase{{n: 2, m: 2}}
	for i := 0; i < 99; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
