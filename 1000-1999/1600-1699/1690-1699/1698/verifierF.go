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
	n int
	a []int
	b []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
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
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1698F.go")
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
	tests = append(tests, Test{n: 1, a: []int{1}, b: []int{1}})
	tests = append(tests, Test{n: 2, a: []int{1, 2}, b: []int{2, 1}})
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(n) + 1
		}
		b := append([]int(nil), a...)
		rng.Shuffle(n, func(i, j int) { b[i], b[j] = b[j], b[i] })
		tests = append(tests, Test{n: n, a: a, b: b})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
