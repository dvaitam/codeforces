package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "refK.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1252K.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
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

type TestCase string

func randAB(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'A'
		} else {
			b[i] = 'B'
		}
	}
	return string(b)
}

func genTests() []TestCase {
	rng := rand.New(rand.NewSource(11))
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		s := randAB(rng, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		sb.WriteString(s)
		sb.WriteByte('\n')
		for j := 0; j < q; j++ {
			if rng.Intn(2) == 0 {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				fmt.Fprintf(&sb, "1 %d %d\n", l, r)
			} else {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				A := rng.Intn(10)
				B := rng.Intn(10)
				fmt.Fprintf(&sb, "2 %d %d %d %d\n", l, r, A, B)
			}
		}
		tests = append(tests, TestCase(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
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
	for i, t := range tests {
		input := string(t)
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
