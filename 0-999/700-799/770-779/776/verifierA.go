package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// generateTest generates the i-th test case input and expected output
func generateTest(i int) (string, string) {
	// deterministic seed for reproducibility
	rand.Seed(int64(i))
	a := fmt.Sprintf("nameA%02d", i)
	b := fmt.Sprintf("nameB%02d", i)
	n := i%5 + 1 // 1..5 days
	var input strings.Builder
	fmt.Fprintf(&input, "%s %s\n%d\n", a, b, n)

	// We'll alternate killings between a and b
	names := []string{a, b}
	var output strings.Builder
	fmt.Fprintf(&output, "%s %s\n", names[0], names[1])

	for j := 0; j < n; j++ {
		killIdx := j % 2
		killName := names[killIdx]
		newName := fmt.Sprintf("rep%02d_%02d", i, j)
		fmt.Fprintf(&input, "%s %s\n", killName, newName)
		names[killIdx] = newName
		fmt.Fprintf(&output, "%s %s\n", names[0], names[1])
	}
	return input.String(), output.String()
}

func runTest(binary string, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	for i := 1; i <= 100; i++ {
		input, expect := generateTest(i)
		got, err := runTest(binary, input)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		// Normalize whitespace for comparison
		expectTrim := strings.TrimSpace(expect)
		if got != expectTrim {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%sGot:\n%s\n", i, input, expectTrim, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
