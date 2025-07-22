package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func solve(s string) string {
	pos := strings.IndexByte(s, '0')
	if pos != -1 {
		s = s[:pos] + s[pos+1:]
	} else {
		s = s[:len(s)-1]
	}
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	s = s[i:]
	if s == "" {
		s = "0"
	}
	return s
}

func generateTests() []testCase {
	var tests []testCase
	for l := 2; len(tests) < 100; l++ {
		start := 1 << (l - 1)
		end := 1<<l - 1
		for num := start; num <= end && len(tests) < 100; num++ {
			s := fmt.Sprintf("%b", num)
			tests = append(tests, testCase{in: s + "\n", out: solve(s)})
		}
	}
	return tests
}

func runTest(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(tc.in))
	}()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	// if binary is go source, build it first
	if strings.HasSuffix(bin, ".go") {
		tmp, err := ioutil.TempFile("", "solbin*")
		if err != nil {
			fmt.Println("cannot create temp file:", err)
			os.Exit(1)
		}
		tmp.Close()
		exec.Command("go", "build", "-o", tmp.Name(), bin).Run()
		bin = tmp.Name()
		defer os.Remove(bin)
	}

	tests := generateTests()
	for i, tc := range tests {
		got, err := runTest(bin, tc)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
