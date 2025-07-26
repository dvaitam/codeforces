package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testF struct {
	a string
	b string
}

func genTestsF() []testF {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testF, 100)
	for i := range tests {
		la := r.Intn(10) + 1
		lb := r.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < la; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte(')')
			}
		}
		a := sb.String()
		sb.Reset()
		for j := 0; j < lb; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte(')')
			}
		}
		b := sb.String()
		tests[i] = testF{a: a, b: b}
	}
	return tests
}

func buildOracle() (string, error) {
	exe := filepath.Join(os.TempDir(), "oracle1272F.bin")
	cmd := exec.Command("go", "build", "-o", exe, "1272F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return exe, nil
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	tests := genTestsF()
	for i, tc := range tests {
		input := tc.a + "\n" + tc.b + "\n"
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
