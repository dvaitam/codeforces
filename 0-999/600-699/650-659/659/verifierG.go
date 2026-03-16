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
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	outBin := filepath.Join(os.TempDir(), "modelG.bin")
	content, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(os.TempDir(), "ref659G.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", outBin, cppPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build model (c++) failed: %v\n%s", err, o)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", outBin, refPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build model failed: %v\n%s", err, o)
		}
	}
	return outBin, nil
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
