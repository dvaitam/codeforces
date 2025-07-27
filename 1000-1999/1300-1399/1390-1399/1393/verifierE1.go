package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testE struct {
	words []string
}

func generateTests() []testE {
	r := rand.New(rand.NewSource(46))
	tests := make([]testE, 100)
	letters := []rune("abc")
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 1
		words := make([]string, n)
		for j := 0; j < n; j++ {
			l := r.Intn(5) + 1
			var sb strings.Builder
			for k := 0; k < l; k++ {
				sb.WriteRune(letters[r.Intn(len(letters))])
			}
			words[j] = sb.String()
		}
		tests[i] = testE{words: words}
	}
	return tests
}

func (t testE) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.words)))
	for _, w := range t.words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", oracle, "1393E1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed\ninput:\n%sexpected: %s got: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
