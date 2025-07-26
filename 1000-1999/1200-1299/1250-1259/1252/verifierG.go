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
	ref := "refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1252G.go")
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

func genTests() []TestCase {
	rng := rand.New(rand.NewSource(7))
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		N := rng.Intn(3) + 1
		M := rng.Intn(3) + 1
		Q := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", N, M, Q)
		for j := 0; j < N; j++ {
			fmt.Fprintf(&sb, "%d ", rng.Intn(10))
		}
		sb.WriteByte('\n')
		for j := 0; j < M; j++ {
			R := rng.Intn(3) + 1
			fmt.Fprintf(&sb, "%d ", R)
			for k := 0; k < R; k++ {
				fmt.Fprintf(&sb, "%d ", rng.Intn(10))
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < Q; j++ {
			X := rng.Intn(M) + 1
			Y := rng.Intn(3)
			Z := rng.Intn(10)
			fmt.Fprintf(&sb, "%d %d %d\n", X, Y, Z)
		}
		tests = append(tests, TestCase(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
