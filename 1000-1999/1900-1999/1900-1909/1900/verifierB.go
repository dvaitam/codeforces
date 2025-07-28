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
	expected string
}

func solveOne(a, b, c int) (int, int, int) {
	one := 0
	two := 0
	three := 0
	if b%2 == c%2 {
		one = 1
	}
	if a%2 == c%2 {
		two = 1
	}
	if a%2 == b%2 {
		three = 1
	}
	return one, two, three
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(10) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		a := rng.Intn(100) + 1
		b := rng.Intn(100) + 1
		c := rng.Intn(100) + 1
		in.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
		o, tw, th := solveOne(a, b, c)
		out.WriteString(fmt.Sprintf("%d %d %d\n", o, tw, th))
	}
	return testCase{input: in.String(), expected: out.String()}
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1\n1 1 1\n", expected: "1 1 1\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
