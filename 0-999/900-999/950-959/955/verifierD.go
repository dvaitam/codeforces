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
	n, m, k int
	s, t    string
}

func (t Test) Input() string {
	return fmt.Sprintf("%d %d %d\n%s\n%s\n", t.n, t.m, t.k, t.s, t.t)
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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "955D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(3))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(n) + 1
		k := rng.Intn(m) + 1
		s := randomString(rng, n)
		t := randomString(rng, m)
		tests = append(tests, Test{n: n, m: m, k: k, s: s, t: t})
	}
	tests = append(tests,
		Test{n: 1, m: 1, k: 1, s: "a", t: "a"},
		Test{n: 5, m: 3, k: 2, s: "abcde", t: "acd"},
		Test{n: 4, m: 4, k: 4, s: "aaaa", t: "bbbb"},
		Test{n: 6, m: 2, k: 1, s: "abcdef", t: "af"},
		Test{n: 8, m: 5, k: 3, s: "abababab", t: "babab"},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:%sexpected:%sgot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
