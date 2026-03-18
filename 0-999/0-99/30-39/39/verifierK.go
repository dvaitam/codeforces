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

var refBin string

func prepareReference() (string, func(), error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		refPath = "39K.go"
	}
	content, err := os.ReadFile(refPath)
	if err != nil {
		return "", nil, fmt.Errorf("cannot read reference at %s: %v", refPath, err)
	}
	if strings.Contains(string(content), "#include") {
		tmpCpp := filepath.Join(os.TempDir(), fmt.Sprintf("ref39K_%d.cpp", time.Now().UnixNano()))
		if err := os.WriteFile(tmpCpp, content, 0644); err != nil {
			return "", nil, err
		}
		tmpBinPath := tmpCpp + ".bin"
		cmd := exec.Command("g++", "-O2", "-o", tmpBinPath, tmpCpp)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmpCpp)
			return "", nil, fmt.Errorf("C++ compile failed: %v\n%s", err, string(out))
		}
		os.Remove(tmpCpp)
		return tmpBinPath, func() { os.Remove(tmpBinPath) }, nil
	}
	tmpGo := filepath.Join(os.TempDir(), fmt.Sprintf("ref39K_%d.go", time.Now().UnixNano()))
	if err := os.WriteFile(tmpGo, content, 0644); err != nil {
		return "", nil, err
	}
	tmpBinPath := tmpGo + ".bin"
	cmd := exec.Command("go", "build", "-o", tmpBinPath, tmpGo)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpGo)
		return "", nil, fmt.Errorf("Go build failed: %v\n%s", err, string(out))
	}
	os.Remove(tmpGo)
	return tmpBinPath, func() { os.Remove(tmpBinPath) }, nil
}

func solveK(input string) string {
	cmd := exec.Command(refBin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

func ok(grid [][]byte, r, c int) bool {
	n := len(grid)
	m := len(grid[0])
	for i := r - 1; i <= r+1; i++ {
		for j := c - 1; j <= c+1; j++ {
			if i >= 0 && i < n && j >= 0 && j < m {
				if grid[i][j] == '*' {
					return false
				}
			}
		}
	}
	return true
}

func generateCaseK(rng *rand.Rand) string {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	k := rng.Intn(2) + 1
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	placed := 0
	attempts := 0
	for placed < k && attempts < 100 {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if grid[r][c] == '.' && ok(grid, r, c) {
			grid[r][c] = '*'
			placed++
		}
		attempts++
	}
	k = placed
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	return sb.String()
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	var cleanup func()
	var err error
	refBin, cleanup, err = prepareReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare reference: %v\n", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseK(rng)
	}
	for i, tc := range cases {
		expect := solveK(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
