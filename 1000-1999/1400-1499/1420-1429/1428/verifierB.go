package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTestsB() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		var sb strings.Builder
		for j := 0; j < n; j++ {
			r := rand.Intn(3)
			if r == 0 {
				sb.WriteByte('<')
			} else if r == 1 {
				sb.WriteByte('>')
			} else {
				sb.WriteByte('-')
			}
		}
		s := sb.String()
		tests[i] = fmt.Sprintf("1\n%d\n%s\n", n, s)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
    bin := os.Args[1]
    tests := genTestsB()
    for i, input := range tests {
        out, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        // Compute expected directly from input
        lines := strings.Split(strings.TrimSpace(input), "\n")
        if len(lines) < 3 {
            fmt.Fprintf(os.Stderr, "malformed test %d\n", i+1)
            os.Exit(1)
        }
        s := strings.TrimSpace(lines[2])
        exp := expected1428B(s)
        got := strings.TrimSpace(out)
        if got != fmt.Sprintf("%d", exp) {
            fmt.Fprintf(os.Stderr, "test %d: expected %d got %s\n", i+1, exp, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

func expected1428B(s string) int {
    n := len(s)
    hasL, hasR := false, false
    for i := 0; i < n; i++ {
        if s[i] == '<' { hasL = true }
        if s[i] == '>' { hasR = true }
    }
    if !hasL || !hasR {
        return n
    }
    ok := make([]bool, n)
    for i := 0; i < n; i++ {
        if s[i] == '-' {
            ok[i] = true
            ok[(i+1)%n] = true
        }
    }
    cnt := 0
    for i := 0; i < n; i++ {
        if ok[i] { cnt++ }
    }
    return cnt
}
