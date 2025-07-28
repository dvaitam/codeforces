package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseC struct {
	n       int
	parents []int
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

func expectedC(n int, parents []int) int {
	cnt := make([]int, n)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		cnt[p-1]++
	}
	arr := []int{}
	for _, c := range cnt {
		if c > 0 {
			arr = append(arr, c)
		}
	}
	arr = append(arr, 1)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	m := len(arr)
	for i := 0; i < m; i++ {
		arr[i] -= m - i
		if arr[i] < 0 {
			arr[i] = 0
		}
	}
	b := []int{}
	for _, v := range arr {
		if v > 0 {
			b = append(b, v)
		}
	}
	if len(b) == 0 {
		return m
	}
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	lo, hi := 0, n+5
	for lo < hi {
		mid := (lo + hi) / 2
		sum := 0
		for _, v := range b {
			if v > mid {
				sum += v - mid
			}
		}
		if sum <= mid {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return m + lo
}

func generateTestsC() ([]testCaseC, []byte) {
	rng := rand.New(rand.NewSource(3))
	t := 100
	tests := make([]testCaseC, t)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 2
		parents := make([]int, n-1)
		for j := 2; j <= n; j++ {
			parents[j-2] = rng.Intn(j-1) + 1
		}
		tests[i] = testCaseC{n: n, parents: parents}
		fmt.Fprintf(&buf, "%d\n", n)
		for j := 0; j < n-1; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", parents[j])
		}
		buf.WriteByte('\n')
	}
	return tests, buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests, input := generateTestsC()
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
		exp := expectedC(tc.n, tc.parents)
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
