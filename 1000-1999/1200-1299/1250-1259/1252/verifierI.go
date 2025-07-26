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
	ref := "refI.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1252I.go")
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
	rng := rand.New(rand.NewSource(9))
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		N := rng.Intn(3)
		xL := 0
		yL := 0
		xR := rng.Intn(5) + 1
		yR := rng.Intn(5) + 1
		xs := rng.Intn(xR-xL+1) + xL
		ys := rng.Intn(yR-yL+1) + yL
		xt := rng.Intn(xR-xL+1) + xL
		yt := rng.Intn(yR-yL+1) + yL
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", N, xL, yL, xR, yR)
		fmt.Fprintf(&sb, "%d %d\n", xs, ys)
		fmt.Fprintf(&sb, "%d %d\n", xt, yt)
		for j := 0; j < N; j++ {
			xi := rng.Intn(xR-xL+1) + xL
			yi := rng.Intn(yR-yL+1) + yL
			ri := rng.Intn(3) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", xi, yi, ri)
		}
		tests = append(tests, TestCase(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
