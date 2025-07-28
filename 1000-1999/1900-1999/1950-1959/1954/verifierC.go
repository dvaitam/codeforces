package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	bin := "refC.bin"
	cmd := exec.Command("go", "build", "-o", bin, "1954C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, out)
	}
	return bin, nil
}

func runBin(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(out), fmt.Errorf("time limit exceeded")
	}
	return string(out), err
}

func generateTests() []string {
	rand.Seed(3)
	tests := make([]string, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		var sb strings.Builder
		var sb2 strings.Builder
		for j := 0; j < n; j++ {
			d := byte(rand.Intn(9) + 1)
			sb.WriteByte('0' + d)
			d2 := byte(rand.Intn(9) + 1)
			sb2.WriteByte('0' + d2)
		}
		tests[i] = fmt.Sprintf("%s\n%s", sb.String(), sb2.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := fmt.Sprintf("%d\n%s\n", len(tests), strings.Join(tests, "\n"))

	expOut, err := runBin("./"+refBin, input)
	if err != nil {
		fmt.Println("reference run error:", err)
		os.Exit(1)
	}
	candOut, err := runBin(cand, input)
	if err != nil {
		fmt.Println("candidate run error:", err)
		os.Exit(1)
	}

	exp := bytes.Fields([]byte(expOut))
	got := bytes.Fields([]byte(candOut))
	if len(exp) != len(got) {
		fmt.Printf("output length mismatch: expected %d got %d\n", len(exp), len(got))
		os.Exit(1)
	}
	for i := range exp {
		if !bytes.Equal(exp[i], got[i]) {
			fmt.Printf("test %d failed\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
