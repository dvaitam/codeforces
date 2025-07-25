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
	expected int
}

func expectedDays(c, v0, v1, a, l int) int {
	days := 1
	read := v0
	speed := v0
	for read < c {
		days++
		speed += a
		if speed > v1 {
			speed = v1
		}
		read += speed - l
	}
	return days
}

func buildCase(c, v0, v1, a, l int) testCase {
	input := fmt.Sprintf("%d %d %d %d %d\n", c, v0, v1, a, l)
	return testCase{input: input, expected: expectedDays(c, v0, v1, a, l)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	c := rng.Intn(1000) + 1
	v0 := rng.Intn(1000) + 1
	l := rng.Intn(v0)
	v1 := v0 + rng.Intn(1000-v0+1)
	a := rng.Intn(1001)
	return buildCase(c, v0, v1, a, l)
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// deterministic edge cases
	cases = append(cases, buildCase(1, 1, 1, 0, 0))
	cases = append(cases, buildCase(5, 4, 4, 0, 0))
	cases = append(cases, buildCase(12, 4, 5, 1, 3))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
