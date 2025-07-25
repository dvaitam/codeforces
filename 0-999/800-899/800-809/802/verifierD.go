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
	exe, err := os.CreateTemp("", "refD*")
	if err != nil {
		return "", err
	}
	exe.Close()
	os.Remove(exe.Name())
	cmd := exec.Command("go", "build", "-o", exe.Name(), "802D.go")
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

func genTests() []string {
	rand.Seed(4)
	tests := make([]string, 0, 102)
	for i := 0; i < 100; i++ {
		vals := make([]int, 250)
		for j := range vals {
			vals[j] = rand.Intn(5)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		for j, v := range vals {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	vals := make([]int, 250)
	tests = append(tests, "1\n"+strings.Repeat("0 ", 249)+"0\n")
	for j := range vals {
		vals[j] = 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	for j, v := range vals {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	tests = append(tests, sb.String())
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
