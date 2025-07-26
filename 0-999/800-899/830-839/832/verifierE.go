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
	n       int
	m       int
	base    []string
	q       int
	queries []string
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, s := range tc.base {
		sb.WriteString(s + "\n")
	}
	sb.WriteString(fmt.Sprintf("%d\n", tc.q))
	for _, b := range tc.queries {
		sb.WriteString(b + "\n")
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
	cmd := exec.Command("go", "build", "-o", ref, "832E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randString(alphabet string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func genTests() []Test {
	rand.Seed(time.Now().UnixNano())
	tests := make([]Test, 0, 100)
	letters := "abcde"
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 1
		m := rand.Intn(4) + 1
		base := make([]string, n)
		for j := 0; j < n; j++ {
			base[j] = randString(letters, m)
		}
		q := rand.Intn(4) + 1
		queries := make([]string, q)
		for j := 0; j < q; j++ {
			queries[j] = randString(letters, m)
		}
		tests = append(tests, Test{n: n, m: m, base: base, q: q, queries: queries})
	}
	// edge
	tests = append(tests, Test{n: 1, m: 1, base: []string{"a"}, q: 1, queries: []string{"a"}})
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
