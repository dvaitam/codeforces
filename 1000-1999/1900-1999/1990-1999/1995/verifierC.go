package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const embeddedTestcasesC = `2 5 6
6 5 5 9 8 3 8
5 3 1 10 3 9
1 6
1 4
6 8 10 4 8 9 3
6 6 3 8 9 1 9
1 9
3 1 2 2
4 10 6 10 8
3 7 9 6
6 2 3 6 1 3 3
1 6
5 4 1 7 1 5
4 1 10 3 6
1 7
1 8
3 10 10 5
6 5 10 8 7 3 1
4 5 4 7 2
3 2 2 1
3 1 3 7
5 1 6 8 9 8
4 2 1 9 7
3 1 9 2
1 6
3 2 8 1
2 9 5
1 1
4 6 3 9 3
2 3 3
6 4 10 6 1 8 7
1 4
2 5 6
2 4 6
2 3 7
4 6 10 3 7
5 1 3 10 1 7
6 3 3 1 1 6 9
1 1
1 2
5 10 3 3 7 1
4 7 10 6 4
2 6 9
2 9 7
1 3
4 10 6 2 7
4 4 8 7 4
4 4 8 7 10
1 5
3 9 6 9
1 10
5 8 4 5 1 10
3 7 2 9
1 3
6 7 1 7 7 7 2
6 8 10 8 3 3 6
4 7 3 10 5
5 2 6 6 3 6
4 9 1 4 5
2 10 6
3 7 1 5
5 7 1 7 5 7
6 4 6 3 3 2 10
3 3 1 7
5 7 8 2 2 7
5 9 3 3 3 4
2 4 1
5 3 8 6 10 5
6 6 2 7 5 3 6
6 9 6 9 3 7 9
3 4 7 6
4 8 9 5 7
4 2 3 3 1
5 10 9 2 4 10
6 10 9 2 5 10 3
4 2 1 1 2
3 8 6 2
6 8 6 10 5 8 4
2 9 9
1 9
2 1 3
5 7 10 4 8 7
3 1 10 3
4 3 8 10 1
4 2 10 7 6
2 9 8
1 8
5 2 5 9 8 9
6 7 8 5 3 4 9
3 3 5 10
2 8 2
1 8
4 10 7 9 2
3 8 4 2
3 3 6 2
2 1 3
5 9 4 1 1 7
5 10 8 2 8 9
3 6 2 1`

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// solve1995CSingle is the core logic from 1995C.go applied to one test case.
func solve1995CSingle(a []int64) int64 {
	n := len(a)
	total := int64(0)
	prevBase := a[0]
	prevK := 0
	for i := 1; i < n; i++ {
		currA := a[i]
		var currK int
		if currA == 1 {
			if prevBase != 1 {
				return -1
			}
			currK = 0
		} else {
			minK := math.MaxInt
			found := false
			start := maxInt(0, prevK-6)
			end := prevK + 6
			for ck := start; ck <= end; ck++ {
				if isGE(currA, ck, prevBase, prevK) {
					if ck < minK {
						minK = ck
					}
					found = true
				}
			}
			if !found {
				return -1
			}
			currK = minK
		}
		total += int64(currK)
		prevBase = currA
		prevK = currK
	}
	return total
}

// solve1995CInput runs the embedded 1995C solution on the provided input.
func solve1995CInput(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return "", fmt.Errorf("read t: %w", err)
	}
	var out strings.Builder
	for tt := 0; tt < t; tt++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return "", fmt.Errorf("read n for test %d: %w", tt+1, err)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &a[i]); err != nil {
				return "", fmt.Errorf("read a[%d] for test %d: %w", i, tt+1, err)
			}
		}
		res := solve1995CSingle(a)
		fmt.Fprintf(&out, "%d\n", res)
	}
	return out.String(), nil
}

func isGE(a int64, k int, pb int64, pk int) bool {
	if k == pk {
		return a >= pb
	}
	if k > pk {
		d := k - pk
		exp := int64(1 << uint(d))
		powa := new(big.Int).Exp(big.NewInt(a), big.NewInt(exp), nil)
		return powa.Cmp(big.NewInt(pb)) >= 0
	}
	d := pk - k
	exp := int64(1 << uint(d))
	powpb := new(big.Int).Exp(big.NewInt(pb), big.NewInt(exp), nil)
	return big.NewInt(a).Cmp(powpb) >= 0
}

func parseTests() ([][]int64, error) {
	var tests [][]int64
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesC))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n: %w", err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("invalid number of values for n=%d", n)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value at pos %d: %w", i+1, err)
			}
			arr[i] = val
		}
		tests = append(tests, arr)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tests [][]int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, arr := range tests {
		fmt.Fprintf(&b, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	expectedOut, err := solve1995CInput(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run embedded 1995C solution: %v\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("runtime error: %v\nstderr: %s\n", err, stderr.String())
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanWords)
	expScan := bufio.NewScanner(strings.NewReader(expectedOut))
	expScan.Split(bufio.ScanWords)
	for i := 0; i < len(tests); i++ {
		if !expScan.Scan() {
			fmt.Printf("missing expected output for test %d\n", i+1)
			os.Exit(1)
		}
		expectStr := expScan.Text()
		expect, err := strconv.ParseInt(expectStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse expected output for test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		gotStr := outScan.Text()
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", i+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
