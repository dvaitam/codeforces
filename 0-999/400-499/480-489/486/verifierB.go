package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeB(A [][]int) [][]int {
	m := len(A)
	n := len(A[0])
	B := make([][]int, m)
	for i := 0; i < m; i++ {
		B[i] = make([]int, n)
		rowOR := 0
		for _, v := range A[i] {
			if v == 1 {
				rowOR = 1
				break
			}
		}
		for j := 0; j < n; j++ {
			colOR := 0
			if rowOR == 1 {
				B[i][j] = 1
				continue
			}
			for k := 0; k < m; k++ {
				if A[k][j] == 1 {
					colOR = 1
					break
				}
			}
			if colOR == 1 {
				B[i][j] = 1
			}
		}
	}
	return B
}

func validAFromB(B [][]int) ([][]int, bool) {
	m := len(B)
	n := len(B[0])
	rowAll := make([]int, m)
	colAll := make([]int, n)
	for i := 0; i < m; i++ {
		all1 := 1
		for j := 0; j < n; j++ {
			if B[i][j] == 0 {
				all1 = 0
				break
			}
		}
		rowAll[i] = all1
	}
	for j := 0; j < n; j++ {
		all1 := 1
		for i := 0; i < m; i++ {
			if B[i][j] == 0 {
				all1 = 0
				break
			}
		}
		colAll[j] = all1
	}
	A := make([][]int, m)
	for i := 0; i < m; i++ {
		A[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if rowAll[i] == 1 && colAll[j] == 1 {
				A[i][j] = 1
			}
		}
	}
	B2 := computeB(A)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if B2[i][j] != B[i][j] {
				return nil, false
			}
		}
	}
	return A, true
}

func parseMatrix(lines []string, m, n int) [][]int {
	A := make([][]int, m)
	for i := 0; i < m; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != n {
			return nil
		}
		A[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if fields[j] == "1" {
				A[i][j] = 1
			} else if fields[j] == "0" {
				A[i][j] = 0
			} else {
				return nil
			}
		}
	}
	return A
}

func verifyCase(bin string, B [][]int) error {
	m := len(B)
	n := len(B[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m, n))
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte(byte('0' + B[i][j]))
		}
		sb.WriteByte('\n')
	}
	got, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	outLines := strings.Split(strings.TrimSpace(got), "\n")
	if len(outLines) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.TrimSpace(outLines[0])
	if strings.ToUpper(first) == "NO" {
		if _, ok := validAFromB(B); ok {
			return fmt.Errorf("solver answered NO but solution exists")
		}
		return nil
	}
	if strings.ToUpper(first) != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if len(outLines)-1 != m {
		return fmt.Errorf("expected %d rows of matrix", m)
	}
	A := parseMatrix(outLines[1:], m, n)
	if A == nil {
		return fmt.Errorf("failed to parse matrix")
	}
	B2 := computeB(A)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if B2[i][j] != B[i][j] {
				return fmt.Errorf("produced matrix A does not match B")
			}
		}
	}
	return nil
}

func genValidCase(rng *rand.Rand, m, n int) [][]int {
	A := make([][]int, m)
	for i := 0; i < m; i++ {
		A[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 1 {
				A[i][j] = 1
			}
		}
	}
	return computeB(A)
}

func genInvalidCase(rng *rand.Rand, m, n int) [][]int {
	// generate random B until it's invalid
	for {
		B := make([][]int, m)
		for i := 0; i < m; i++ {
			B[i] = make([]int, n)
			for j := 0; j < n; j++ {
				B[i][j] = rng.Intn(2)
			}
		}
		if _, ok := validAFromB(B); !ok {
			return B
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		m := rng.Intn(4) + 1 // 1..4
		n := rng.Intn(4) + 1
		var B [][]int
		if rng.Float64() < 0.7 {
			B = genValidCase(rng, m, n)
		} else {
			B = genInvalidCase(rng, m, n)
		}
		if err := verifyCase(bin, B); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
