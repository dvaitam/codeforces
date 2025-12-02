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

type op struct{ x, y int }

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

func solveOps(A [][]int) ([]op, bool) {
	n := len(A)
	m := len(A[0])
	B := make([][]int, n)
	for i := range B {
		B[i] = make([]int, m)
	}
	var ops []op
	for i := 0; i+1 < n; i++ {
		for j := 0; j+1 < m; j++ {
			if A[i][j] == 1 && A[i][j+1] == 1 && A[i+1][j] == 1 && A[i+1][j+1] == 1 {
				ops = append(ops, op{i + 1, j + 1})
				B[i][j], B[i][j+1], B[i+1][j], B[i+1][j+1] = 1, 1, 1, 1
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] != B[i][j] {
				return nil, false
			}
		}
	}
	return ops, true
}

func applyOps(n, m int, ops []op) [][]int {
	B := make([][]int, n)
	for i := range B {
		B[i] = make([]int, m)
	}
	for _, o := range ops {
		x := o.x - 1
		y := o.y - 1
		if x < 0 || x+1 >= n || y < 0 || y+1 >= m {
			continue
		}
		B[x][y] = 1
		B[x][y+1] = 1
		B[x+1][y] = 1
		B[x+1][y+1] = 1
	}
	return B
}

func equal(A, B [][]int) bool {
	for i := range A {
		for j := range A[i] {
			if A[i][j] != B[i][j] {
				return false
			}
		}
	}
	return true
}

func genTest() [][]int {
	n := 2 + rand.Intn(10)
	m := 2 + rand.Intn(10)
	A := make([][]int, n)
	for i := range A {
		A[i] = make([]int, m)
		for j := range A[i] {
			if rand.Intn(2) == 0 {
				A[i][j] = 1
			}
		}
	}
	return A
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx := 1; idx <= 100; idx++ {
		A := genTest()
		n := len(A)
		m := len(A[0])
		_, ok := solveOps(A)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				input.WriteString(fmt.Sprintf("%d", A[i][j]))
				if j+1 < m {
					input.WriteByte(' ')
				}
			}
			input.WriteByte('\n')
		}
		cmd := exec.Command(cand)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("test %d: no output\n", idx)
			os.Exit(1)
		}
		first := outScan.Text()
		if first == "-1" {
			if ok {
				fmt.Printf("test %d: expected solution but got -1\n", idx)
				os.Exit(1)
			}
			if outScan.Scan() {
				fmt.Printf("test %d: extra output after -1\n", idx)
				os.Exit(1)
			}
			continue
		}
		if !ok {
			fmt.Printf("test %d: expected -1 but got solution\n", idx)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(first)
		if k > n*m {
			fmt.Printf("test %d: too many operations %d > %d\n", idx, k, n*m)
			os.Exit(1)
		}
		var candOps []op
		for i := 0; i < k; i++ {
			if !outScan.Scan() {
				fmt.Printf("test %d: missing op x\n", idx)
				os.Exit(1)
			}
			x, _ := strconv.Atoi(outScan.Text())
			if !outScan.Scan() {
				fmt.Printf("test %d: missing op y\n", idx)
				os.Exit(1)
			}
			y, _ := strconv.Atoi(outScan.Text())
			candOps = append(candOps, op{x, y})
		}
		if outScan.Scan() {
			fmt.Printf("test %d: extra output\n", idx)
			os.Exit(1)
		}
		B := applyOps(n, m, candOps)
		if !equal(A, B) {
			fmt.Printf("test %d: wrong operations\n", idx)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
