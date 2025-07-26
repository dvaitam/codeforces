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
	exe, err := os.CreateTemp("", "refK*")
	if err != nil {
		return "", err
	}
	exe.Close()
	os.Remove(exe.Name())
	cmd := exec.Command("go", "build", "-o", exe.Name(), "802K.go")
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

func randTree(n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, rand.Intn(n)+1))
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		w := rand.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", p, i, w))
	}
	return sb.String()
}

func genTests() []string {
	rand.Seed(11)
	tests := make([]string, 0, 102)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		tests = append(tests, randTree(n))
	}
	tests = append(tests, randTree(1))
	tests = append(tests, randTree(2))
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:%sexpected: %s got: %s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
