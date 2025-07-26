package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	words []string
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", len(t.words)))
	for i, w := range t.words {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildRef() (string, error) {
	ref := "refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1259D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randWord(r *rand.Rand) string {
	n := r.Intn(6) + 1
	b := make([]byte, n)
	for i := range b {
		if r.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func genTests() []Test {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := r.Intn(6) + 1
		words := make([]string, 0, n)
		m := make(map[string]bool)
		for len(words) < n {
			w := randWord(r)
			if !m[w] {
				m[w] = true
				words = append(words, w)
			}
		}
		tests = append(tests, Test{words})
	}
	tests = append(tests, Test{[]string{"0"}})
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
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
