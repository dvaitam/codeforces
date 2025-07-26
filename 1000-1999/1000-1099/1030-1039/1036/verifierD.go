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
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1036D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d", rand.Intn(10)+1))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d", rand.Intn(10)+1))
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests,
		Test{"1\n5\n1\n5\n"},
		Test{"2\n1 2\n1\n3\n"},
		Test{"3\n1 1 1\n3\n1 1 1\n"},
		Test{"4\n1 2 3 4\n4\n4 3 2 1\n"},
		Test{"2\n10 20\n3\n5 10 15\n"},
	)
	return tests
}

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
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
