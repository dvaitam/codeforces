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

type testCase struct {
	input  string
	expect int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(30)
	letters := make([]string, n)
	seen := make(map[rune]bool)
	for i := 0; i < n; i++ {
		r := rune('a' + rng.Intn(26))
		letters[i] = string(r)
		seen[r] = true
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, l := range letters {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(l)
	}
	b.WriteByte('}')
	b.WriteByte('\n')
	return testCase{input: b.String(), expect: len(seen)}
}

func runCase(bin string, tc testCase) error {
	out, err := run(bin, tc.input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, got)
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
	cases := []testCase{
		{"{}\n", 0},
		{"{a}\n", 1},
		{"{a, b, c}\n", 3},
		{"{a, a, b, b, c}\n", 3},
	}
	for i := 0; i < 120; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
