package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "204E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

type testCase struct {
	n    int
	k    int
	strs []string
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(5))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{{n: 1, k: 1, strs: []string{"a"}}}
	tests = append(tests, fixed...)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for len(tests) < 100 {
		n := rng.Intn(4) + 1
		k := rng.Intn(n) + 1
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			l := rng.Intn(5) + 1
			b := make([]byte, l)
			for j := 0; j < l; j++ {
				b[j] = letters[rng.Intn(26)]
			}
			strs[i] = string(b)
		}
		tests = append(tests, testCase{n: n, k: k, strs: strs})
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := generateTests()
	for i, t := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
		for _, s := range t.strs {
			input.WriteString(fmt.Sprintf("%s\n", s))
		}
		expOut, err := run(oracle, input.String())
		if err != nil {
			fmt.Printf("oracle error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expOut {
			fmt.Printf("test %d failed. Expected %q got %q\n", i+1, expOut, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
