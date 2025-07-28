package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func digitSum(x int) int {
	sum := 0
	for x > 0 {
		sum += x % 10
		x /= 10
	}
	return sum
}

func solve(x, k int) int {
	for {
		if digitSum(x)%k == 0 {
			return x
		}
		x++
	}
}

func generateTests() [][2]int {
	rng := rand.New(rand.NewSource(1))
	tests := make([][2]int, 0, 100)
	// some fixed edge cases
	tests = append(tests, [2]int{1, 1}, [2]int{1, 10}, [2]int{1000000000, 10})
	for len(tests) < 100 {
		x := rng.Intn(1000000000) + 1
		k := rng.Intn(10) + 1
		tests = append(tests, [2]int{x, k})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candA")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc[0], tc[1])
		expected := strconv.Itoa(solve(tc[0], tc[1]))
		output, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if output != expected {
			fmt.Printf("case %d failed: x=%d k=%d expected %s got %s\n", i+1, tc[0], tc[1], expected, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
