package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCmd(exe, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	refBin := "./refA_bin"
	cmd := exec.Command("go", "build", "-o", refBin, "612A.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refBin, nil
}

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		p := r.Intn(n) + 1
		q := r.Intn(n) + 1
		sb := strings.Builder{}
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(len(letters))])
		}
		s := sb.String()
		tests[i] = fmt.Sprintf("%d %d %d\n%s\n", n, p, q, s)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		expect, err := runCmd(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
