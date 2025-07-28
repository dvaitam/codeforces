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

const mod int64 = 998244353

type Test struct {
	n int
	s int
	a []int
	b []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.s))
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range t.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1698E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref failed: %v: %s", err, out)
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 102)
	// deterministic small case
	tests = append(tests, buildCase(3, 1, []int{2, 1, 3}))
	tests = append(tests, buildCase(3, 2, []int{2, 1, 3}))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		s := rng.Intn(n) + 1
		perm := rng.Perm(n)
		a := make([]int, n)
		for j := range a {
			a[j] = perm[j] + 1
		}
		// build b
		b := make([]int, n)
		for j := range b {
			b[j] = -1
		}
		numbers := rng.Perm(n)
		k := rng.Intn(n + 1)
		usedIdx := rng.Perm(n)[:k]
		for idx, pos := range usedIdx {
			b[pos] = numbers[idx] + 1
		}
		tests = append(tests, Test{n: n, s: s, a: a, b: b})
	}
	return tests
}

func buildCase(n int, s int, a []int) Test {
	b := make([]int, n)
	for i := range b {
		b[i] = -1
	}
	return Test{n: n, s: s, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
