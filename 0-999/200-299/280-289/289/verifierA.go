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

type segment struct {
	l, r int
}

func expectedAnswer(n int, k int, segs []segment) int64 {
	var sum int64
	for _, s := range segs {
		sum += int64(s.r - s.l + 1)
	}
	rem := sum % int64(k)
	if rem == 0 {
		return 0
	}
	return int64(k) - rem
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(10) + 1
	k := rng.Intn(100) + 1
	segs := make([]segment, n)
	start := rng.Intn(200) - 100
	for i := 0; i < n; i++ {
		length := rng.Intn(20) + 1
		end := start + length - 1
		segs[i] = segment{l: start, r: end}
		start = end + rng.Intn(20) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d\n", s.l, s.r)
	}
	expected := expectedAnswer(n, k, segs)
	return sb.String(), expected
}

func runCase(exe string, input string, expected int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("cannot parse output: %v\n%s", err, outStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// a couple of deterministic cases
	cases := []struct {
		input    string
		expected int64
	}{
		{"1 5\n0 0\n", 5 - (1 % 5)},
		{"2 3\n0 0\n5 7\n", (3 - (4 % 3))},
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, struct {
			input    string
			expected int64
		}{in, exp})
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
