package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildReference() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "27A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
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

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %q", out)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val <= 0 {
		return 0, fmt.Errorf("answer must be positive, got %d", val)
	}
	return val, nil
}

func deterministicCases() []string {
	cases := []string{
		"1\n1\n",
		"1\n2\n",
		"3\n2 3 4\n",
		"5\n2 3 4 5 6\n",
	}
	var sb strings.Builder
	sb.WriteString("3000\n")
	for i := 1; i <= 3000; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	sb.WriteByte('\n')
	cases = append(cases, sb.String())
	return cases
}

func generateRandomCase(rng *rand.Rand) string {
	n := rng.Intn(3000) + 1
	perm := rng.Perm(3000)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(perm[i] + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func verifyCase(candidate, ref, input string) error {
	refOut, err := runProgram(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseAnswer(refOut)
	if err != nil {
		return fmt.Errorf("bad reference output: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		return err
	}
	got, err := parseAnswer(candOut)
	if err != nil {
		return err
	}
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
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicCases()
	for i := 0; i < 200; i++ {
		tests = append(tests, generateRandomCase(rng))
	}

	for idx, input := range tests {
		if err := verifyCase(candidate, ref, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
