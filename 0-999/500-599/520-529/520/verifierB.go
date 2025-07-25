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

type testCaseB struct {
	n int
	m int
}

func generateCase(rng *rand.Rand) (string, testCaseB) {
	n := rng.Intn(10000) + 1
	m := rng.Intn(10000) + 1
	input := fmt.Sprintf("%d %d\n", n, m)
	return input, testCaseB{n: n, m: m}
}

func expected(tc testCaseB) int {
	n, m := tc.n, tc.m
	if n >= m {
		return n - m
	}
	steps := 0
	for m > n {
		if m%2 == 0 {
			m /= 2
		} else {
			m++
		}
		steps++
	}
	return steps + n - m
}

func runCase(bin string, input string, tc testCaseB) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
