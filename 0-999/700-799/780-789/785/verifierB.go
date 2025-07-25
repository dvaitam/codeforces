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

type Interval struct {
	l, r int64
}

type Test struct {
	n, m int
	a, b []Interval
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.n) + "\n")
	for _, in := range t.a {
		sb.WriteString(fmt.Sprintf("%d %d\n", in.l, in.r))
	}
	sb.WriteString(strconv.Itoa(t.m) + "\n")
	for _, in := range t.b {
		sb.WriteString(fmt.Sprintf("%d %d\n", in.l, in.r))
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
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "785B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 103)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		a := make([]Interval, n)
		b := make([]Interval, m)
		for j := 0; j < n; j++ {
			l := rand.Int63n(50)
			r := l + rand.Int63n(10)
			a[j] = Interval{l, r}
		}
		for j := 0; j < m; j++ {
			l := rand.Int63n(50)
			r := l + rand.Int63n(10)
			b[j] = Interval{l, r}
		}
		tests = append(tests, Test{n: n, m: m, a: a, b: b})
	}
	tests = append(tests, Test{n: 1, m: 1, a: []Interval{{1, 2}}, b: []Interval{{3, 4}}})
	tests = append(tests, Test{n: 1, m: 1, a: []Interval{{1, 5}}, b: []Interval{{2, 3}}})
	tests = append(tests, Test{n: 1, m: 1, a: []Interval{{1, 3}}, b: []Interval{{1, 3}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
