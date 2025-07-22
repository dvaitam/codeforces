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

func runCandidate(bin, input string) (string, error) {
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

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		a := rng.Intn(179) + 1
		sb.WriteString(fmt.Sprintf("%d\n", a))
		if 360%(180-a) == 0 {
			expected[i] = "YES"
		} else {
			expected[i] = "NO"
		}
	}
	return sb.String(), expected
}

func checkOutput(out string, expected []string) error {
	fields := strings.Fields(out)
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d lines got %d", len(expected), len(fields))
	}
	for i, exp := range expected {
		if fields[i] != exp {
			return fmt.Errorf("line %d: expected %s got %s", i+1, exp, fields[i])
		}
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

	// add sample test case
	sample := "3\n30\n60\n90\n"
	sampleExp := []string{"NO", "YES", "YES"}
	out, err := runCandidate(bin, sample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sample failed: %v\n", err)
		os.Exit(1)
	}
	if err := checkOutput(out, sampleExp); err != nil {
		fmt.Fprintf(os.Stderr, "sample failed: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkOutput(out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
