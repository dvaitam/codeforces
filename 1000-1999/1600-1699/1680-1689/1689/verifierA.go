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

type Test struct {
	n, m, k int
	a, b    string
}

func (t Test) Input() string {
	return fmt.Sprintf("1\n%d %d %d\n%s\n%s\n", t.n, t.m, t.k, t.a, t.b)
}

func buildOracle() (string, error) {
	ref := "oracleA"
	cmd := exec.Command("go", "build", "-o", ref, "1689A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	lettersA := []rune("abcdefghijklm")
	lettersB := []rune("nopqrstuvwxyz")
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		k := rng.Intn(10) + 1
		a := make([]rune, n)
		b := make([]rune, m)
		for j := 0; j < n; j++ {
			a[j] = lettersA[rng.Intn(len(lettersA))]
		}
		for j := 0; j < m; j++ {
			b[j] = lettersB[rng.Intn(len(lettersB))]
		}
		tests[i] = Test{n, m, k, string(a), string(b)}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
