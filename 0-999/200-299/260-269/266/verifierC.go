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

func parseCase(line string) (int, [][2]int, string) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return 0, nil, ""
	}
	n, _ := strconv.Atoi(fields[0])
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
	return n, pairs, sb.String()
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
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, pairs, input := parseCase(line)
		// build initial matrix
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
		}
		for _, p := range pairs {
			mat[p[0]-1][p[1]-1] = 1
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing m\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(scan.Text())
		if err != nil || m < 0 || m > 100000 {
			fmt.Fprintf(os.Stderr, "case %d: invalid m\n", idx)
			os.Exit(1)
		}
		ops := make([][3]int, m)
		for i := 0; i < m; i++ {
			for j := 0; j < 3; j++ {
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d: incomplete op\n", idx)
					os.Exit(1)
				}
				v, err := strconv.Atoi(scan.Text())
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d: invalid op value\n", idx)
					os.Exit(1)
				}
				ops[i][j] = v
			}
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx)
			os.Exit(1)
		}
		if err := apply(mat, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if !valid(mat) {
			fmt.Fprintf(os.Stderr, "case %d: resulting matrix invalid\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
