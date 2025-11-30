package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one per line).
const embeddedTestcases = `8 1 2 5 4 6 8 5 1 2 6 1 3 3 5
8 3 8 7 6 8 6 3 6 1 6 6 3 1 3
8 5 8 6 8 4 6 8 3 4 8 8 6 7 8
8 5 4 5 7 1 4 2 6 8 6 8 5 3 5
7 2 7 6 1 5 6 7 5 1 3 5 2
4 3 2 1 3 1 4
8 1 5 3 1 6 8 1 4 1 7 1 6 2 5
5 4 5 2 1 1 3 2 4
7 7 1 5 7 1 4 6 3 5 3 2 5
3 2 3 1 3
6 6 2 1 2 1 4 2 6 1 3
6 1 5 6 1 6 4 2 5 1 3
8 2 1 6 5 4 2 7 3 1 7 7 5 1 3
3 2 3 3 2
3 3 1 3 2
5 2 1 5 1 1 4 4 3
2 1 2
4 3 2 2 1 1 4
2 2 1
4 3 1 2 4 2 1
6 3 4 3 1 6 4 2 3 4 1
2 1 2
7 7 1 1 5 4 3 6 1 6 4 4 2
4 3 4 4 2 4 3
6 6 5 1 5 4 6 3 2 5 2
8 7 4 2 1 6 8 6 4 2 3 7 2 3 5
7 2 4 4 6 7 3 1 7 5 6 5 3
6 2 4 4 3 4 6 6 4 5 1
4 3 1 1 3 2 1
2 2 1
7 6 5 1 4 4 2 2 3 7 5 5 2
5 4 2 1 5 3 5 5 2
8 1 2 5 8 6 8 1 8 2 6 6 3 2 8
8 3 8 1 5 3 1 8 1 5 3 7 5 7 8
4 4 2 3 4 1 4
7 2 4 6 5 7 3 6 7 6 3 1 3
4 2 4 4 1 2 1
4 3 1 1 2 4 2
2 1 2
8 3 8 8 4 6 4 8 3 2 6 8 6 8 5
4 3 2 1 2 4 2
8 3 4 5 7 7 6 2 6 5 6 7 5 6 3
4 2 3 4 1 3 4
3 1 2 1 3
8 7 1 4 6 6 4 2 3 8 3 4 8 7 5
4 2 4 1 3 2 1
4 3 1 1 4 4 3
3 2 3 3 1
4 3 1 1 4 4 2
6 5 4 4 6 6 4 4 2 4 5
3 3 1 1 3
2 2 1
4 1 2 1 3 4 1
4 1 4 4 1 4 2
3 3 1 1 3
6 1 2 3 4 6 5 5 1 5 3
2 2 1
8 7 4 4 3 1 4 7 3 2 3 7 6 5 6
5 2 4 5 4 1 4 1 5
4 2 4 1 2 2 1
3 3 1 1 2
7 5 1 6 4 6 7 2 6 3 6 1 6
8 4 3 6 1 5 4 4 2 5 2 3 6 2 8
5 1 3 1 2 4 1 4 3
5 4 5 5 4 2 4 3 4
4 2 4 4 2 4 3
4 3 2 1 2 3 4
6 3 4 3 1 2 3 5 3 1 3
2 1 2
6 5 4 4 6 1 6 2 5 5 2
4 4 2 2 3 1 4
8 7 4 8 7 3 7 5 4 6 7 7 2 4 1
5 3 2 5 1 4 2 1 5
6 3 1 4 6 4 5 1 6 3 5
7 5 1 2 3 2 6 3 2 4 1 5 2
7 6 2 2 1 3 7 4 6 1 6 3 5
7 3 4 1 5 3 1 3 7 5 7 3 6
4 1 3 2 1 4 3
5 1 3 5 1 4 2 1 5
8 5 1 1 4 4 2 1 7 2 6 6 3 4 7
2 2 1
2 2 1
5 3 1 5 4 4 5 5 2
8 8 4 1 5 6 1 5 4 8 3 7 5 6 3
8 3 8 7 1 8 4 1 8 4 8 8 6 7 5
4 1 2 4 1 3 4
5 5 3 5 4 4 1 4 2
5 3 1 5 4 1 2 4 5
5 2 3 1 3 3 4 3 5
8 7 4 2 8 1 5 6 7 8 3 1 6 8 2
5 5 3 4 5 2 3 5 1
4 3 2 1 2 4 3
4 3 2 2 1 2 4
5 4 5 5 2 1 2 4 3
3 3 1 1 3
2 2 1
7 7 1 3 4 1 5 2 3 4 5 1 3
5 2 4 1 3 5 1 4 2
8 1 3 2 7 5 8 6 8 4 6 1 4 4 1
7 7 1 3 7 5 7 1 7 7 2 1 3`

type op struct{ t, i, j int }

// Embedded solver logic from 266C.go for reference.
func solveOps(n int, pairs [][2]int) []op {
	aX := make([]int, n)
	aY := make([]int, n)
	for i, p := range pairs {
		aX[i+1] = p[0]
		aY[i+1] = p[1]
	}
	var ops []op
	for i := 1; i < n; i++ {
		if aX[i] != i+1 {
			old := aX[i]
			newRow := i + 1
			for j := i + 1; j < n; j++ {
				if aX[j] == old {
					aX[j] = newRow
				} else if aX[j] == newRow {
					aX[j] = old
				}
			}
			ops = append(ops, op{1, old, newRow})
			aX[i] = newRow
		}
		if aY[i] > i {
			old := aY[i]
			newCol := i
			for j := i + 1; j < n; j++ {
				if aY[j] == old {
					aY[j] = newCol
				} else if aY[j] == newCol {
					aY[j] = old
				}
			}
			ops = append(ops, op{2, old, newCol})
			aY[i] = newCol
		}
	}
	return ops
}

// Candidate runner.
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

func parseCase(line string) (int, [][2]int, string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return 0, nil, "", fmt.Errorf("empty line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, "", err
	}
	if len(fields) != 1+2*(n-1) {
		return 0, nil, "", fmt.Errorf("length mismatch")
	}
	pairs := make([][2]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		x, _ := strconv.Atoi(fields[1+2*i])
		y, _ := strconv.Atoi(fields[2+2*i])
		pairs = append(pairs, [2]int{x, y})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return n, pairs, sb.String(), nil
}

func swapRows(mat [][]int, i, j int) {
	mat[i], mat[j] = mat[j], mat[i]
}

func swapCols(mat [][]int, i, j int) {
	for k := range mat {
		mat[k][i], mat[k][j] = mat[k][j], mat[k][i]
	}
}

func apply(mat [][]int, ops [][3]int) error {
	n := len(mat)
	for _, op := range ops {
		t, i, j := op[0], op[1]-1, op[2]-1
		if i < 0 || i >= n || j < 0 || j >= n || i == j {
			return fmt.Errorf("invalid indices")
		}
		if t == 1 {
			swapRows(mat, i, j)
		} else if t == 2 {
			swapCols(mat, i, j)
		} else {
			return fmt.Errorf("invalid op type")
		}
	}
	return nil
}

func valid(mat [][]int) bool {
	count := 0
	n := len(mat)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] == 1 {
				if i <= j {
					return false
				}
				count++
			}
		}
	}
	return count == n-1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		n, pairs, input, err := parseCase(strings.TrimSpace(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
		}
		for _, p := range pairs {
			mat[p[0]-1][p[1]-1] = 1
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing m\n", idx+1)
			os.Exit(1)
		}
		m, err := strconv.Atoi(scan.Text())
		if err != nil || m < 0 || m > 100000 {
			fmt.Fprintf(os.Stderr, "case %d: invalid m\n", idx+1)
			os.Exit(1)
		}
		ops := make([][3]int, m)
		for i := 0; i < m; i++ {
			for j := 0; j < 3; j++ {
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d: incomplete op\n", idx+1)
					os.Exit(1)
				}
				v, err := strconv.Atoi(scan.Text())
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d: invalid op value\n", idx+1)
					os.Exit(1)
				}
				ops[i][j] = v
			}
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		if err := apply(mat, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !valid(mat) {
			fmt.Fprintf(os.Stderr, "case %d: resulting matrix invalid\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
