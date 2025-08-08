package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func parseInput(input []byte) ([]int64, error) {
	r := bytes.NewReader(input)
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return nil, err
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(r, &a[i]); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func bestValue(a []int64) (int64, []int64) {
	n := len(a)
	p := make([]int64, n+1)
	for i := 0; i < n; i++ {
		p[i+1] = p[i] + a[i]
	}
	prefVal := make([]int64, n+1)
	prefVal[0] = p[0]
	for j := 1; j <= n; j++ {
		if p[j] > prefVal[j-1] {
			prefVal[j] = p[j]
		} else {
			prefVal[j] = prefVal[j-1]
		}
	}
	suffVal := make([]int64, n+1)
	suffVal[n] = p[n]
	for j := n - 1; j >= 0; j-- {
		if p[j] >= suffVal[j+1] {
			suffVal[j] = p[j]
		} else {
			suffVal[j] = suffVal[j+1]
		}
	}
	best := int64(-1 << 63)
	for j := 0; j <= n; j++ {
		cur := prefVal[j] - p[j] + suffVal[j]
		if cur > best {
			best = cur
		}
	}
	return best, p
}

func runTests(dir, binary string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		a, err := parseInput(input)
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}
		best, prefix := bestValue(a)
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}
		var d0, d1, d2 int
		if _, err := fmt.Fscan(strings.NewReader(string(out)), &d0, &d1, &d2); err != nil {
			return fmt.Errorf("%s: could not parse output %q", filepath.Base(inFile), strings.TrimSpace(string(out)))
		}
		n := len(a)
		if d0 < 0 || d0 > d1 || d1 > d2 || d2 > n {
			return fmt.Errorf("%s: invalid indices %d %d %d", filepath.Base(inFile), d0, d1, d2)
		}
		val := prefix[d0] - prefix[d1] + prefix[d2]
		if val != best {
			return fmt.Errorf("%s: expected value %d but got %d", filepath.Base(inFile), best, val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "C")
	if err := runTests(testDir, binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
