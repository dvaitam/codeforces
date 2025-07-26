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
	words []string
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.words)))
	for _, w := range t.words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return sb.String()
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
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "965E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randWord() string {
	l := rand.Intn(6) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func genTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 103)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		set := make(map[string]struct{})
		words := make([]string, 0, n)
		for len(words) < n {
			w := randWord()
			if _, ok := set[w]; ok {
				continue
			}
			set[w] = struct{}{}
			words = append(words, w)
		}
		tests = append(tests, Test{words})
	}
	tests = append(tests,
		Test{[]string{"a"}},
		Test{[]string{"abc", "abcd", "ab", "abcdef"}},
		Test{[]string{"x", "y", "z"}},
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
		input := tc.Input()
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
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
