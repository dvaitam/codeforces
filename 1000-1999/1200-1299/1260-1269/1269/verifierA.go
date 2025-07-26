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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isComposite(x int) bool {
	if x < 4 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return true
		}
	}
	return false
}

type testCase int

func generateTests(rng *rand.Rand) []testCase {
	tests := []testCase{1, 2, 3, 4, 5, 10, 17, 37, 9999999, 10000000}
	for len(tests) < 100 {
		tests = append(tests, testCase(rng.Intn(10000000)+1))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)

	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected two numbers got %q\ninput:\n%s", i+1, out, input)
			os.Exit(1)
		}
		var a, b int
		if _, err := fmt.Sscan(fields[0], &a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad number %q\ninput:\n%s", i+1, fields[0], input)
			os.Exit(1)
		}
		if _, err := fmt.Sscan(fields[1], &b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad number %q\ninput:\n%s", i+1, fields[1], input)
			os.Exit(1)
		}
		if a-b != int(n) {
			fmt.Fprintf(os.Stderr, "case %d failed: a-b != n (%d-%d=%d)\ninput:\n%s", i+1, a, b, a-b, input)
			os.Exit(1)
		}
		if a < 2 || a > 1000000000 || b < 2 || b > 1000000000 {
			fmt.Fprintf(os.Stderr, "case %d failed: values out of range a=%d b=%d\ninput:\n%s", i+1, a, b, input)
			os.Exit(1)
		}
		if !isComposite(a) || !isComposite(b) {
			fmt.Fprintf(os.Stderr, "case %d failed: numbers must be composite a=%d b=%d\ninput:\n%s", i+1, a, b, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
