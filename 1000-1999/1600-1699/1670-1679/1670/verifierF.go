package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1670F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []string {
	rand.Seed(6)
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		l := rand.Int63n(50) + 1
		r := l + rand.Int63n(50)
		z := rand.Int63n(64)
		tests = append(tests, fmt.Sprintf("1\n%d %d %d %d\n", n, l, r, z))
	}
	tests = append(tests, "1\n1 1 1 0\n")
	tests = append(tests, "1\n2 1 3 1\n")
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		exp, err := run(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
