package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCmd(exe, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", exe)
	} else {
		cmd = exec.CommandContext(ctx, exe)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	refBin := "./refA_bin"
	cmd := exec.Command("go", "build", "-o", refBin, "612A.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refBin, nil
}

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		p := r.Intn(n) + 1
		q := r.Intn(n) + 1
		sb := strings.Builder{}
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(len(letters))])
		}
		s := sb.String()
		tests[i] = fmt.Sprintf("%d %d %d\n%s\n", n, p, q, s)
	}
	return tests
}

func checkSolution(input, output string) error {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	// If the solution claims -1, it is structurally valid (logic in main handles if it SHOULD be -1)
	if fields[0] == "-1" {
		return nil
	}

	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}

	if len(fields) != k+1 {
		return fmt.Errorf("expected %d substrings, got %d", k, len(fields)-1)
	}

	// Parse input
	inFields := strings.Fields(input)
	n, _ := strconv.Atoi(inFields[0])
	p, _ := strconv.Atoi(inFields[1])
	q, _ := strconv.Atoi(inFields[2])
	s := inFields[3]

	reconstructed := ""
	for i := 0; i < k; i++ {
		sub := fields[i+1]
		if len(sub) != p && len(sub) != q {
			return fmt.Errorf("substring %d (%q) has length %d, expected %d or %d", i+1, sub, len(sub), p, q)
		}
		reconstructed += sub
	}

	if reconstructed != s {
		return fmt.Errorf("reconstructed string %q does not match input %q", reconstructed, s)
	}
	if len(reconstructed) != n {
		return fmt.Errorf("total length %d != n %d", len(reconstructed), n)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	// Add the specific failing case to the front of the test list
	tests = append([]string{"6 6 3\nINvNSQ\n"}, tests...)

	for idx, input := range tests {
		expect, err := runCmd(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		expectTrim := strings.TrimSpace(expect)
		gotTrim := strings.TrimSpace(got)

		if expectTrim == "-1" {
			if gotTrim != "-1" {
				fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: -1\ngot:\n%s\n", idx+1, input, got)
				os.Exit(1)
			}
		} else {
			if gotTrim == "-1" {
				fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected solution\ngot: -1\n", idx+1, input)
				os.Exit(1)
			}
			if err := checkSolution(input, got); err != nil {
				fmt.Printf("wrong answer on test %d: %v\ninput:\n%s\ngot:\n%s\n", idx+1, err, input, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
