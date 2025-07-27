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
	input string
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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1458D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randString(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			b.WriteByte('0')
		} else {
			b.WriteByte('1')
		}
	}
	return b.String()
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		t := rand.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			lenS := rand.Intn(10) + 1
			sb.WriteString(randString(lenS) + "\n")
		}
		tests = append(tests, Test{sb.String()})
	}
	// edge
	tests = append(tests, Test{"1\n0\n"})
	return tests
}

func normalize(s string) string { return strings.TrimSpace(s) }

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if normalize(exp) != normalize(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
