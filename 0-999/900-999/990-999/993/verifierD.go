package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func compileReference() (string, error) {
	refPath := filepath.Join(os.TempDir(), fmt.Sprintf("refD_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", refPath, "993D.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(1))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)+1))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := compileReference()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		return
	}
	defer os.Remove(ref)
	tests := generateTests()
	for i, t := range tests {
		exp, err := runBinary(ref, t)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(candidate, t)
		if err != nil {
			fmt.Printf("tested binary failed on test %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}
