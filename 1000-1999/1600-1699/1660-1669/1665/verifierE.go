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

type query struct{ l, r int }

type testCaseE struct {
	n       int
	arr     []int
	q       int
	queries []query
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

func expectedE(arr []int, l, r int) int {
	ans := int(^uint(0) >> 1)
	for i := l; i < r; i++ {
		for j := i + 1; j <= r; j++ {
			val := arr[i] | arr[j]
			if val < ans {
				ans = val
			}
		}
	}
	return ans
}

func generateTestsE() ([]testCaseE, []byte) {
	rng := rand.New(rand.NewSource(5))
	t := 100
	tests := make([]testCaseE, t)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(15) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(1 << 10)
		}
		q := rng.Intn(10) + 1
		qs := make([]query, q)
		fmt.Fprintf(&buf, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", arr[j])
		}
		buf.WriteByte('\n')
		fmt.Fprintf(&buf, "%d\n", q)
		for j := 0; j < q; j++ {
			l := rng.Intn(n - 1)
			r := l + 1 + rng.Intn(n-l-1)
			qs[j] = query{l: l + 1, r: r + 1}
			fmt.Fprintf(&buf, "%d %d\n", l+1, r+1)
		}
		tests[i] = testCaseE{n: n, arr: arr, q: q, queries: qs}
	}
	return tests, buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests, input := generateTestsE()
	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanWords)
	for idx, tc := range tests {
		for j := 0; j < tc.q; j++ {
			if !scanner.Scan() {
				fmt.Printf("missing output for test %d query %d\n", idx+1, j+1)
				os.Exit(1)
			}
			got, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Printf("invalid integer on test %d query %d: %v\n", idx+1, j+1, err)
				os.Exit(1)
			}
			exp := expectedE(tc.arr, tc.queries[j].l-1, tc.queries[j].r-1)
			if got != exp {
				fmt.Printf("test %d query %d failed: expected %d got %d\n", idx+1, j+1, exp, got)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
