package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func compileRef() (string, error) {
	exe, err := os.CreateTemp("", "refC*")
	if err != nil {
		return "", err
	}
	exe.Close()
	os.Remove(exe.Name())
	cmd := exec.Command("go", "build", "-o", exe.Name(), "802C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile reference: %v\n%s", err, string(out))
	}
	return exe.Name(), nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func joinInts(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genTests() []string {
	rand.Seed(3)
	tests := make([]string, 0, 102)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(20) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(n) + 1
		}
		c := make([]int, n)
		for j := range c {
			c[j] = rand.Intn(100)
		}
		input := fmt.Sprintf("%d %d\n%s\n%s\n", n, k, joinInts(a), joinInts(c))
		tests = append(tests, input)
	}
	tests = append(tests, "1 1\n1\n0\n")
	tests = append(tests, "2 1\n1 2\n5 5\n")
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("reference compile failed:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, in := range tests {
		exp, err := runProg(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
