package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const numTests = 100

func buildRef(dir string) (string, error) {
	ref := filepath.Join(dir, "refC_bin")
	cmd := exec.Command("go", "build", "-o", ref, "896C.go")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	ref, err := buildRef(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for t := 1; t <= numTests; t++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		seed := rand.Int63n(1000000)
		vmax := rand.Int63n(1000) + 1
		input := fmt.Sprintf("%d %d %d %d\n", n, m, seed, vmax)
		inBytes := []byte(input)

		out, err := runBinary(target, inBytes)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		refOut, err := runBinary(ref, inBytes)
		if err != nil {
			fmt.Printf("Test %d: reference runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(refOut) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", t, input, refOut, out)
			os.Exit(1)
		}
	}
	fmt.Printf("Passed %d tests\n", numTests)
}
