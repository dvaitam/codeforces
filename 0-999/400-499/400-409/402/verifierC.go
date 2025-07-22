package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierC.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refC")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "402C.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open(filepath.Join(dir, "testcasesC.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	idx := 0
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		idx++
		input := fmt.Sprintf("1\n%s\n", line)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", idx, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: reference error: %v\n", idx, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", idx, input, refOut, candOut)
			os.Exit(1)
		}
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
