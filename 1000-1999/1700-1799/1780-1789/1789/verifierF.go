package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(20) + 1 // 1..20
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rand.Intn(26)))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func buildReference() (string, func(), error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", nil, fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", nil, fmt.Errorf("cannot read reference source: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "1789F-ref")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	binPath := filepath.Join(tmpDir, "ref_1789F")

	if strings.Contains(string(content), "#include") {
		// C++ source saved with .go extension
		cppPath := filepath.Join(tmpDir, "ref.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			cleanup()
			return "", nil, err
		}
		cmd := exec.Command("g++", "-O2", "-o", binPath, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("g++ build failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", binPath, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
		}
	}

	return binPath, cleanup, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]

	ref, cleanup, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
