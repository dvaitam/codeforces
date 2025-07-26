package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n, k, m, d int64
}

func (t Test) Input() string {
	return fmt.Sprintf("%d %d %d %d\n", t.n, t.k, t.m, t.d)
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "965C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 103)
	for i := 0; i < 100; i++ {
		k := int64(rand.Intn(1000) + 2)
		m := int64(rand.Intn(1000) + 1)
		d := int64(rand.Intn(100) + 1)
		n := int64(rand.Intn(1000000) + 1)
		maxN := k * m * d
		if maxN > 0 && n > maxN {
			n = maxN
		}
		if n < 2 {
			n = 2
		}
		if k > n {
			k = n
		}
		if m > n {
			m = n
		}
		if d > n {
			d = n
		}
		tests = append(tests, Test{n, k, m, d})
	}
	tests = append(tests,
		Test{2, 2, 1, 1},
		Test{1000000, 2, 1000000, 1},
		Test{1000000, 1000000, 1000000, 1000},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
