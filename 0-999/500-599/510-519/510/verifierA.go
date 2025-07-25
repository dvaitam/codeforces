package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
)

func runTests(dir, binary string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		outFile := inFile[:len(inFile)-3] + ".out"
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		expected, err := os.ReadFile(outFile)
		if err != nil {
			return err
		}
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}
		if string(out) != string(expected) {
			return fmt.Errorf("%s: expected\n%sbut got\n%s", filepath.Base(inFile), expected, out)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "A")
	if err := runTests(testDir, binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
