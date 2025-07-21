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

func solveCase(lines []string) string {
	n := len(lines) / 2
	m := make(map[string]struct{})
	for i := 0; i < n; i++ {
		species := lines[2*i]
		color := lines[2*i+1]
		key := species + " " + color
		m[key] = struct{}{}
	}
	return fmt.Sprintf("%d", len(m))
}

var species = []string{"oak", "maple", "birch", "pine", "elm"}
var colors = []string{"red", "green", "yellow", "orange", "brown"}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	lines := make([]string, 0, 2*n)
	for i := 0; i < n; i++ {
		sp := species[rng.Intn(len(species))]
		co := colors[rng.Intn(len(colors))]
		sb.WriteString(sp + " " + co + "\n")
		lines = append(lines, sp, co)
	}
	return testCase{input: sb.String(), expected: solveCase(lines)}
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
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
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
	var cases []testCase
	// deterministic simple case
	cases = append(cases, testCase{input: "1\noak red\n", expected: "1"})
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
