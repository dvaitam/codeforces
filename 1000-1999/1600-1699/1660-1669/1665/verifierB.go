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

type testCaseB struct {
	n   int
	arr []int
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

func expectedB(arr []int) int {
	freq := make(map[int]int)
	mx := 0
	n := len(arr)
	for _, v := range arr {
		freq[v]++
		if freq[v] > mx {
			mx = freq[v]
		}
	}
	ops := 0
	for mx < n {
		ops++
		diff := n - mx
		add := mx
		if diff < add {
			add = diff
		}
		ops += add
		mx += add
	}
	return ops
}

func generateTestsB() ([]testCaseB, []byte) {
	rng := rand.New(rand.NewSource(2))
	t := 100
	tests := make([]testCaseB, t)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(7) - 3
		}
		tests[i] = testCaseB{n: n, arr: arr}
		fmt.Fprintf(&buf, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", arr[j])
		}
		buf.WriteByte('\n')
	}
	return tests, buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests, input := generateTestsB()
	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanWords)
	for idx, tc := range tests {
		if !scanner.Scan() {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		got, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("invalid integer on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expectedB(tc.arr)
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
