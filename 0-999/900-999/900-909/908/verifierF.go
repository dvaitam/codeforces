package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const solutionF = "908F.go"

// fallback deterministic tests when no stdin provided
func generateTests() []string {
	rand.Seed(42)
	var tests []string
	for i := 0; i < 5; i++ {
		n := rand.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		pos := 0
		for j := 0; j < n; j++ {
			pos += rand.Intn(5) + 1
			color := []byte{'R', 'B', 'G'}[rand.Intn(3)]
			sb.WriteString(fmt.Sprintf("%d %c\n", pos, color))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runCmd(cmd []string, input string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func runCandidate(path string, input string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		return runCmd([]string{"go", "run", path}, input)
	}
	os.Chmod(path, 0755)
	return runCmd([]string{path}, input)
}

func expectedOutput(input string) (string, error) {
	return runCmd([]string{"go", "run", solutionF}, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]

	// If stdin has content, use it as a single test. Otherwise, run generated tests one by one.
	if data, err := io.ReadAll(os.Stdin); err == nil && len(data) > 0 {
		exp, err := expectedOutput(string(data))
		if err != nil {
			fmt.Printf("reference failed: %v\n", err)
			os.Exit(1)
		}
		got, err := runCandidate(cand, string(data))
		if err != nil {
			fmt.Printf("candidate runtime error: %v\n", err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("mismatch\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", string(data), exp, got)
			os.Exit(1)
		}
		fmt.Println("all tests passed")
		return
	}

	for idx, tc := range generateTests() {
		exp, err := expectedOutput(tc)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(cand, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("mismatch on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
