package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1500C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func applySorts(a [][]int, cols []int) [][]int {
	for i := len(cols) - 1; i >= 0; i-- {
		c := cols[i]
		sort.SliceStable(a, func(x, y int) bool { return a[x][c] < a[y][c] })
	}
	return a
}

func equal(a, b [][]int) bool {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		return false
	}
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for idx := 0; idx < t; idx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		A := make([][]int, n)
		B := make([][]int, n)
		for i := 0; i < n; i++ {
			A[i] = make([]int, m)
			for j := 0; j < m; j++ {
				scan.Scan()
				A[i][j], _ = strconv.Atoi(scan.Text())
			}
		}
		for i := 0; i < n; i++ {
			B[i] = make([]int, m)
			for j := 0; j < m; j++ {
				scan.Scan()
				B[i][j], _ = strconv.Atoi(scan.Text())
			}
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprintf(&input, "%d", A[i][j])
			}
			input.WriteByte('\n')
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprintf(&input, "%d", B[i][j])
			}
			input.WriteByte('\n')
		}
		expectStr, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
			os.Exit(1)
		}
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expFields := strings.Fields(expectStr)
		if expFields[0] == "-1" {
			if strings.TrimSpace(gotStr) != "-1" {
				fmt.Printf("case %d failed: expected -1 got %s\n", idx+1, gotStr)
				os.Exit(1)
			}
			continue
		}
		// expected possible
		gotFields := strings.Fields(gotStr)
		if len(gotFields) < 1 {
			fmt.Printf("case %d: empty output\n", idx+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(gotFields[0])
		if err != nil || len(gotFields) != k+1 {
			fmt.Printf("case %d: invalid output\n", idx+1)
			os.Exit(1)
		}
		cols := make([]int, k)
		for i := 0; i < k; i++ {
			v, _ := strconv.Atoi(gotFields[i+1])
			if v < 1 || v > m {
				fmt.Printf("case %d: column out of range\n", idx+1)
				os.Exit(1)
			}
			cols[i] = v - 1
		}
		AA := make([][]int, n)
		for i := 0; i < n; i++ {
			AA[i] = append([]int(nil), A[i]...)
		}
		AA = applySorts(AA, cols)
		if !equal(AA, B) {
			fmt.Printf("case %d failed: sequence does not transform A to B\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
