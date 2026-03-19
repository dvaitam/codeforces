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

func buildReference() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		src = "1702G1.go"
	}
	// Detect C++ by checking content for #include
	content, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	ref := filepath.Join(os.TempDir(), "refG1.bin")
	if strings.Contains(string(content), "#include") {
		cppFile := filepath.Join(os.TempDir(), "refG1.cpp")
		if err := os.WriteFile(cppFile, content, 0644); err != nil {
			return "", fmt.Errorf("write cpp: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", ref, cppFile)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build c++ ref: %v\n%s", err, out)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", ref, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build go ref: %v\n%s", err, out)
		}
	}
	return ref, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	var input bytes.Buffer
	n := 10
	fmt.Fprintln(&input, n)
	for i := 0; i < n-1; i++ {
		fmt.Fprintf(&input, "%d %d\n", i+1, i+2)
	}
	q := 100
	fmt.Fprintln(&input, q)
	for i := 0; i < q; i++ {
		k := rand.Intn(n) + 1
		fmt.Fprint(&input, k)
		for j := 0; j < k; j++ {
			fmt.Fprintf(&input, " %d", rand.Intn(n)+1)
		}
		fmt.Fprintln(&input)
	}

	refCmd := exec.Command(ref)
	refCmd.Stdin = bytes.NewReader(input.Bytes())
	expOut, err := refCmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed to run reference solution:", err)
		os.Exit(1)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	expLines := strings.Fields(string(expOut))
	gotLines := strings.Fields(string(out))
	if len(expLines) != len(gotLines) {
		fmt.Println("output line count mismatch")
		os.Exit(1)
	}
	for i := range expLines {
		if expLines[i] != gotLines[i] {
			fmt.Printf("test line %d failed: expected %s got %s\n", i+1, expLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
