package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test string

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
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "976A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(0)
	tests := []Test{
		Test("1\n0\n"),
		Test("1\n1\n"),
		Test("4\n1001\n"),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		var sb strings.Builder
		if n == 1 {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		} else {
			sb.WriteByte('1')
			for j := 1; j < n; j++ {
				if rand.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
		}
		tests = append(tests, Test(fmt.Sprintf("%d\n%s\n", n, sb.String())))
	}
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
	for i, t := range tests {
		input := string(t)
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
			fmt.Printf("Test %d failed\nInput:\n%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
