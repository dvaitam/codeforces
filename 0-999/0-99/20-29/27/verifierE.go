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

const maxAnswer = uint64(1e18)

func buildReference() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "27E.go")
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

func parseAnswer(out string) (uint64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %q", out)
	}
	val, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val == 0 || val > maxAnswer {
		return 0, fmt.Errorf("answer %d out of valid range", val)
	}
	return val, nil
}

func deterministicCases() []string {
	values := []int{1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 18, 24, 30, 36, 48, 60, 72, 84, 96, 128, 192, 256, 384, 512, 720, 840, 960, 1000}
	tests := make([]string, 0, len(values))
	for _, v := range values {
		tests = append(tests, fmt.Sprintf("%d\n", v))
	}
	return tests
}

func randomCases(rng *rand.Rand, count int) []string {
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(1000) + 1
		tests = append(tests, fmt.Sprintf("%d\n", n))
	}
	return tests
}

func verifyCase(candidate, ref, input string) error {
	refOut, err := runProgram(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseAnswer(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	tests = append(tests, randomCases(rng, 500)...)

	for idx, input := range tests {
		if err := verifyCase(candidate, ref, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
