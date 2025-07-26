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
	refPath := filepath.Join(os.TempDir(), fmt.Sprintf("refA_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", refPath, "993A.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func generateTests() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		x := rand.Intn(40) - 20
		y := rand.Intn(40) - 20
		size := rand.Intn(10) + 1
		line1 := fmt.Sprintf("%d %d %d %d %d %d %d %d", x, y, x+size, y, x+size, y+size, x, y+size)
		x2 := rand.Intn(40) - 20
		y2 := rand.Intn(40) - 20
		r := rand.Intn(10) + 1
		line2 := fmt.Sprintf("%d %d %d %d %d %d %d %d", x2-r, y2, x2, y2-r, x2+r, y2, x2, y2+r)
		tests[i] = line1 + "\n" + line2 + "\n"
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(candidate, t)
		if err != nil {
			fmt.Printf("tested binary failed on test %d: %v\n", i+1, err)
			return
		}
		if strings.ToLower(strings.TrimSpace(out)) != strings.ToLower(strings.TrimSpace(exp)) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}
