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

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// normalizeTokens extracts all whitespace-separated tokens and joins them with single spaces.
func normalizeTokens(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		refSrc = "773D.go"
	}
	refBin := "./773D_ref"

	// Check if reference is C++ by reading content
	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		// C++ source saved as .go, copy to .cpp and compile
		cppSrc := "./773D_ref.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", fmt.Errorf("write cpp: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", refBin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("compile C++: %v\n%s", err, out)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", refBin, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("compile Go: %v\n%s", err, out)
		}
	}
	return refBin, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				w := rand.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d ", w))
			}
			if i < n-2 {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", t+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", t+1, ":", err)
			os.Exit(1)
		}
		if normalizeTokens(expect) != normalizeTokens(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", t+1, normalizeTokens(expect), normalizeTokens(got))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
