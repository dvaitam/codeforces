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

type testCase struct{ a, b, c, d int }

func solveA(a, b, c, d int) string {
	for t := b; t <= 100000; t += a {
		if t >= d && (t-d)%c == 0 {
			return fmt.Sprintf("%d\n", t)
		}
	}
	return "-1\n"
}

func runCase(bin string, input string, expected string) error {
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
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
	tests := []testCase{
		{1, 1, 1, 1},
		{1, 2, 2, 1},
		{3, 5, 7, 2},
		{5, 20, 3, 30},
		{2, 2, 2, 3},
	}
	for i := 0; i < 100; i++ {
		t := testCase{rng.Intn(100) + 1, rng.Intn(100) + 1, rng.Intn(100) + 1, rng.Intn(100) + 1}
		tests = append(tests, t)
	}
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n%d %d\n", tc.a, tc.b, tc.c, tc.d)
		expected := solveA(tc.a, tc.b, tc.c, tc.d)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
