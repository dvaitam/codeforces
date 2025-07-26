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
	input    string
	expected int64
}

func solveA(n int, x, y int64, s string) int64 {
	zeros := 0
	i := 0
	for i < n {
		if s[i] == '0' {
			zeros++
			for i < n && s[i] == '0' {
				i++
			}
		} else {
			i++
		}
	}
	if zeros == 0 {
		return 0
	}
	cost1 := int64(zeros) * y
	cost2 := y + int64(zeros-1)*x
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	x := rng.Int63n(10)
	y := rng.Int63n(10)
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	expected := solveA(n, x, y, s)
	input := fmt.Sprintf("%d %d %d\n%s\n", n, x, y, s)
	return testCase{input: input, expected: expected}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
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
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
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
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
