package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Query struct {
	l, r int
}

type Test struct {
	n, q int
	qs   []Query
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.q))
	for _, qu := range t.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.l, qu.r))
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
	cmd := exec.Command("go", "build", "-o", ref, "785E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 103)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		q := rand.Intn(20) + 1
		qs := make([]Query, q)
		for j := 0; j < q; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n) + 1
			qs[j] = Query{l: l, r: r}
		}
		tests = append(tests, Test{n: n, q: q, qs: qs})
	}
	tests = append(tests, Test{n: 1, q: 1, qs: []Query{{1, 1}}})
	tests = append(tests, Test{n: 2, q: 1, qs: []Query{{1, 2}}})
	tests = append(tests, Test{n: 3, q: 2, qs: []Query{{1, 2}, {2, 3}}})
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
