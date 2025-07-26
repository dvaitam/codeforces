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
	n int
	s string
}

func (t Test) Input() string {
	return fmt.Sprintf("%d\n%s\n", t.n, t.s)
}

func buildRef() (string, error) {
	ref := "refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1285A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
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
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rng.Intn(1000) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				b[j] = 'L'
			} else {
				b[j] = 'R'
			}
		}
		tests = append(tests, Test{n, string(b)})
	}
	n := 100000
	b := make([]byte, n)
	for j := 0; j < n; j++ {
		if j%2 == 0 {
			b[j] = 'L'
		} else {
			b[j] = 'R'
		}
	}
	tests = append(tests, Test{n, string(b)})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
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
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
