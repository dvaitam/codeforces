package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/source-or-binary")
		os.Exit(1)
	}

	tests := buildTests()

	candidate := os.Args[1]
	binPath, cleanup, err := buildIfNeeded(candidate)
	if err != nil {
		fmt.Printf("build failed: %v\n", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}

	runTests(binPath, tests)
}

func runTests(binary string, tests []struct{ n, m, a int64 }) {
	for i, t := range tests {
		fmt.Printf("running test %d/%d\n", i+1, len(tests))
		expected := tilesNeeded(t.n, t.m, t.a)
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d %d %d\n", t.n, t.m, t.a))
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr:\n%s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		var got int64
		outStr := strings.TrimSpace(outBuf.String())
		fmt.Sscan(outStr, &got)
		if got != expected {
			fmt.Printf("Test %d failed: expected %d got %s\n", i+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func tilesNeeded(n, m, a int64) int64 {
	rows := (n + a - 1) / a
	cols := (m + a - 1) / a
	return rows * cols
}

func buildTests() []struct{ n, m, a int64 } {
	cases := make([]struct{ n, m, a int64 }, 0, 120)
	seed := []struct{ n, m, a int64 }{
		{6, 6, 4}, // sample
		{1, 1, 1}, // trivial
		{1, 2, 3}, // a larger than side
		{1_000_000_000, 1, 1_000_000_000},
		{1_000_000_000, 1_000_000_000, 1_000_000_000},
		{999_999_937, 999_999_929, 2},
		{100, 25, 7},
		{25, 100, 7},
		{99999999, 1234567, 89},
		{33, 44, 5},
		{44, 33, 5},
		{100000, 99999, 17},
	}
	cases = append(cases, seed...)
	for i := int64(0); len(cases) < 110; i++ {
		n := 1 + (i*37)%1_000_000_000
		m := 1 + (i*91)%1_000_000_000
		a := 1 + (i*53)%999_999_900
		if a == 0 {
			a = 1
		}
		cases = append(cases, struct{ n, m, a int64 }{n, m, a})
	}
	return cases
}

func buildIfNeeded(path string) (string, func(), error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		tmp := filepath.Join(os.TempDir(), "candidateA_go.bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", nil, err
		}
		return tmp, func() { _ = os.Remove(tmp) }, nil
	case ".cpp", ".cc", ".cxx", ".c++":
		tmp := filepath.Join(os.TempDir(), "candidateA_cpp.bin")
		cmd := exec.Command("g++", "-std=c++17", "-O2", "-pipe", "-static", "-s", path, "-o", tmp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", nil, err
		}
		return tmp, func() { _ = os.Remove(tmp) }, nil
	case ".py":
		// Run with system python
		return path, nil, nil
	default:
		// Assume already a runnable binary
		return path, nil, nil
	}
}
