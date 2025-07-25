package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildModel() (string, error) {
	modelPath := filepath.Join(filepath.Dir(os.Args[0]), "659G.go")
	out := filepath.Join(os.TempDir(), "modelG.bin")
	cmd := exec.Command("go", "build", "-o", out, modelPath)
	err := cmd.Run()
	return out, err
}

func solveWithModel(model, input string) (string, error) {
	cmd := exec.Command(model)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []string {
	rand.Seed(48)
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		var b strings.Builder
		fmt.Fprintf(&b, "%d\n", n)
		for i := 0; i < n; i++ {
			val := rand.Intn(10) + 1
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
		tests[t] = b.String()
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	model, err := buildModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build model: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(model)
	tests := generateTests()
	for i, t := range tests {
		expect, err := solveWithModel(model, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "model run error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect = strings.TrimSpace(expect)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
