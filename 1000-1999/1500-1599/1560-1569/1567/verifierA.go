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
	n        int
	s        string
	expected string
}

func solveCase(n int, s string) string {
	b := []byte(s)
	for i, c := range b {
		if c == 'U' {
			b[i] = 'D'
		} else if c == 'D' {
			b[i] = 'U'
		}
	}
	return string(b)
}

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

func genTest(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	letters := "LRUD"
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	s := string(b)
	return testCase{n: n, s: s, expected: solveCase(n, s)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n%s\n", tc.n, tc.s))
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
