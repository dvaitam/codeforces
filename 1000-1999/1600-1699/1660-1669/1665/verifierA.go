package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseA struct {
	n int64
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runBinary(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcm(a, b int64) int64 {
	g := gcd(a, b)
	if g == 0 {
		return 0
	}
	return a / g * b
}

func generateTestsA() ([]testCaseA, []byte) {
	rng := rand.New(rand.NewSource(1))
	t := 100
	tests := make([]testCaseA, t)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000_000-4) + 4
		tests[i] = testCaseA{n: n}
		fmt.Fprintf(&buf, "%d\n", n)
	}
	return tests, buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binPath, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests, input := generateTestsA()
	out, err := runBinary(binPath, input)
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanWords)

	for idx, tc := range tests {
		var vals [4]int64
		for i := 0; i < 4; i++ {
			if !scanner.Scan() {
				fmt.Printf("missing value at test %d\n", idx+1)
				os.Exit(1)
			}
			v, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Printf("invalid integer at test %d: %v\n", idx+1, err)
				os.Exit(1)
			}
			vals[i] = v
		}
		a, b, c, d := vals[0], vals[1], vals[2], vals[3]
		if a <= 0 || b <= 0 || c <= 0 || d <= 0 {
			fmt.Printf("test %d has non-positive value\n", idx+1)
			os.Exit(1)
		}
		if a+b+c+d != tc.n {
			fmt.Printf("test %d sum mismatch: got %d expected %d\n", idx+1, a+b+c+d, tc.n)
			os.Exit(1)
		}
		if gcd(a, b) != lcm(c, d) {
			fmt.Printf("test %d gcd/lcm mismatch\n", idx+1)
			os.Exit(1)
		}
	}

	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
