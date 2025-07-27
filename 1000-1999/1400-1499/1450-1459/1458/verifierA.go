package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1458A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 120)
	for i := 0; i < 110; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rand.Intn(100) + 1
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rand.Intn(100) + 1
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	// edge cases
	tests = append(tests, Test{"1 1\n5\n7\n"})
	tests = append(tests, Test{"3 2\n1 1 1\n0 0\n"})
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
