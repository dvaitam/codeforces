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

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1615C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 20)
	for i := 0; i < 19; i++ {
		n := rand.Intn(8) + 1
		var a, b strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				a.WriteByte('0')
			} else {
				a.WriteByte('1')
			}
			if rand.Intn(2) == 0 {
				b.WriteByte('0')
			} else {
				b.WriteByte('1')
			}
		}
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, a.String(), b.String())
		tests = append(tests, Test{input})
	}
	tests = append(tests, Test{"1\n1\n0\n1\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		want, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
