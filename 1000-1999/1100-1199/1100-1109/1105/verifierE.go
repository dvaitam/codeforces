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
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1105E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randName(r *rand.Rand, idx int) string {
	return fmt.Sprintf("f%d", idx)
}

func genTests() []Test {
	r := rand.New(rand.NewSource(4))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		m := r.Intn(5) + 1
		n := r.Intn(20) + 1
		names := make([]string, m)
		for j := 0; j < m; j++ {
			names[j] = randName(r, j)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		sb.WriteString("1\n")
		visited := make([]bool, m)
		for j := 1; j < n; j++ {
			if r.Intn(2) == 0 {
				sb.WriteString("1\n")
			} else {
				idx := r.Intn(m)
				visited[idx] = true
				sb.WriteString("2 " + names[idx] + "\n")
			}
		}
		for idx, v := range visited {
			if !v {
				sb.WriteString("2 " + names[idx] + "\n")
			}
		}
		tests = append(tests, Test{sb.String()})
	}
	// fixed tests
	tests = append(tests,
		Test{"1 1\n1\n"},
		Test{"3 2\n1\n2 a\n2 b\n"},
		Test{"4 1\n1\n2 x\n1\n2 x\n"},
		Test{"5 2\n1\n2 a\n2 b\n1\n2 a\n"},
		Test{"6 3\n1\n2 a\n2 b\n2 c\n1\n2 b\n"},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
