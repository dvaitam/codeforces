package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// expectedAnswer returns the minimum n for the given sharpness x,
// using the same formula as the reference solution.
func expectedAnswer(x int) int {
	if x == 3 {
		return 5
	}
	for i := 0; i < 100; i++ {
		odd := 2*i + 1
		acc := (odd/2)*odd + (odd/2) + 1
		if x <= acc {
			return odd
		}
	}
	return -1
}

func runBinary(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Test all 100 possible inputs exhaustively.
	for x := 1; x <= 100; x++ {
		tc := fmt.Sprintf("%d\n", x)
		expected := fmt.Sprintf("%d", expectedAnswer(x))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "x=%d failed: %v\ninput:\n%s", x, err, tc)
			os.Exit(1)
		}
		if expected != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "x=%d wrong answer\nexpected: %s\ngot: %s\n", x, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
