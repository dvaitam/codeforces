package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refPath = "2000-2999/2100-2199/2130-2139/2138/2138E1.go"

type testCase struct {
	x int
}

func main() {
	if len(os.Args) != 2 {
		if len(os.Args) == 3 && os.Args[1] == "--" {
			os.Args = []string{os.Args[0], os.Args[2]}
		} else {
			fmt.Println("usage: go run verifierE1.go /path/to/binary")
			os.Exit(1)
		}
	}
	bin := os.Args[1]

	tests := buildTests()
	input := renderInput(tests)

	if _, err := runBinary(refPath, input); err != nil {
		fmt.Printf("reference runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Printf("runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	if err := verifyAll(out, tests); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}

func buildTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase
	tests = append(tests, testCase{0}, testCase{1}, testCase{2}, testCase{3}, testCase{5}, testCase{10}, testCase{42}, testCase{10000000})
	for i := 0; i < 40; i++ {
		tests = append(tests, testCase{rand.Intn(10000001)})
	}
	return tests
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.x)
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verifyAll(out string, tests []testCase) error {
	tokens := strings.Fields(out)
	idx := 0
	for caseIdx, tc := range tests {
		if idx >= len(tokens) {
			return fmt.Errorf("case %d: output ended early", caseIdx+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		idx++
		if err != nil {
			return fmt.Errorf("case %d: invalid n: %v", caseIdx+1, err)
		}
		if n < 1 || n > 80 {
			return fmt.Errorf("case %d: n out of range %d", caseIdx+1, n)
		}
		need := n * n
		if idx+need > len(tokens) {
			return fmt.Errorf("case %d: not enough matrix entries", caseIdx+1)
		}
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
			for j := 0; j < n; j++ {
				v, err := strconv.Atoi(tokens[idx])
				idx++
				if err != nil {
					return fmt.Errorf("case %d: invalid entry at (%d,%d): %v", caseIdx+1, i+1, j+1, err)
				}
				if v < -1 || v > 1 {
					return fmt.Errorf("case %d: entry (%d,%d)=%d not in [-1,0,1]", caseIdx+1, i+1, j+1, v)
				}
				mat[i][j] = v
			}
		}
		if err := checkSparsity(mat, caseIdx); err != nil {
			return err
		}
		det := determinant(mat)
		if det.Cmp(big.NewInt(int64(tc.x))) != 0 {
			return fmt.Errorf("case %d: determinant mismatch, got %s expected %d", caseIdx+1, det.String(), tc.x)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extraneous tokens after parsing output (used %d of %d)", idx, len(tokens))
	}
	return nil
}

func checkSparsity(mat [][]int, caseIdx int) error {
	n := len(mat)
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < n; j++ {
			if mat[i][j] != 0 {
				cnt++
			}
		}
		if cnt > 3 {
			return fmt.Errorf("case %d: row %d has %d non-zero entries", caseIdx+1, i+1, cnt)
		}
	}
	for j := 0; j < n; j++ {
		cnt := 0
		for i := 0; i < n; i++ {
			if mat[i][j] != 0 {
				cnt++
			}
		}
		if cnt > 3 {
			return fmt.Errorf("case %d: column %d has %d non-zero entries", caseIdx+1, j+1, cnt)
		}
	}
	return nil
}

func determinant(mat [][]int) *big.Int {
	n := len(mat)
	if n == 1 {
		return big.NewInt(int64(mat[0][0]))
	}
	a := make([][]big.Int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]big.Int, n)
		for j := 0; j < n; j++ {
			a[i][j].SetInt64(int64(mat[i][j]))
		}
	}
	prev := big.NewInt(1)
	sign := int64(1)
	for k := 0; k < n-1; k++ {
		pivot := &a[k][k]
		if pivot.Sign() == 0 {
			swap := -1
			for r := k + 1; r < n; r++ {
				if a[r][k].Sign() != 0 {
					swap = r
					break
				}
			}
			if swap == -1 {
				return big.NewInt(0)
			}
			a[k], a[swap] = a[swap], a[k]
			pivot = &a[k][k]
			sign = -sign
		}
		for i := k + 1; i < n; i++ {
			for j := k + 1; j < n; j++ {
				var num1, num2, res big.Int
				num1.Mul(pivot, &a[i][j])
				num2.Mul(&a[i][k], &a[k][j])
				num1.Sub(&num1, &num2)
				res.Quo(&num1, prev)
				a[i][j].Set(&res)
			}
		}
		prev = new(big.Int).Set(pivot)
	}
	det := new(big.Int).Set(&a[n-1][n-1])
	if sign == -1 {
		det.Neg(det)
	}
	return det
}
