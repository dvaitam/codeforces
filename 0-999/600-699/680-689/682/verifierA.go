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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func solve(n, m int64) int64 {
	cntN := [5]int64{}
	cntM := [5]int64{}
	cntN[0] = n / 5
	cntM[0] = m / 5
	for i := int64(1); i < 5; i++ {
		cntN[i] = n / 5
		if n%5 >= i {
			cntN[i]++
		}
		cntM[i] = m / 5
		if m%5 >= i {
			cntM[i]++
		}
	}
	var ans int64
	for r := 0; r < 5; r++ {
		comp := (5 - r) % 5
		ans += cntN[r] * cntM[comp]
	}
	return ans
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	var got int64
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := solve(tc.n, tc.m)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{1, 1})
	cases = append(cases, testCase{5, 5})
	for i := 0; i < 100; i++ {
		n := rng.Int63n(1_000_000) + 1
		m := rng.Int63n(1_000_000) + 1
		cases = append(cases, testCase{n, m})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d\n", i+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
