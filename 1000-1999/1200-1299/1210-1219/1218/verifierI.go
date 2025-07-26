package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct{ input string }

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refI.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1218I.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(9)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 2
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				if rand.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				if rand.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
		c1 := rand.Intn(n)
		for k := 0; k < n; k++ {
			if k == c1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierI.go /path/to/binary")
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
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
