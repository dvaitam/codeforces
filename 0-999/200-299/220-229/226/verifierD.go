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

func parseMatrix(scanner *bufio.Scanner) (int, int, [][]int, bool) {
	if !scanner.Scan() {
		return 0, 0, nil, false
	}
	header := strings.TrimSpace(scanner.Text())
	if header == "" {
		if !scanner.Scan() {
			return 0, 0, nil, false
		}
		header = strings.TrimSpace(scanner.Text())
	}
	parts := strings.Fields(header)
	if len(parts) < 2 {
		return 0, 0, nil, false
	}
	n, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return 0, 0, nil, false
		}
		rowParts := strings.Fields(scanner.Text())
		row := make([]int, m)
		for j := 0; j < m; j++ {
			v, _ := strconv.Atoi(rowParts[j])
			row[j] = v
		}
		mat[i] = row
	}
	scanner.Scan() // consume blank line
	return n, m, mat, true
}

func runCase(bin string, n, m int, mat [][]int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", mat[i][j]))
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("output should contain two lines")
	}
	rowParts := strings.Fields(lines[0])
	colParts := strings.Fields(lines[1])
	var rCount int
	var rows []int
	if len(rowParts) > 0 {
		rCount, _ = strconv.Atoi(rowParts[0])
		for i := 0; i < rCount; i++ {
			v, _ := strconv.Atoi(rowParts[i+1])
			rows = append(rows, v)
		}
	}
	var cCount int
	var cols []int
	if len(colParts) > 0 {
		cCount, _ = strconv.Atoi(colParts[0])
		for i := 0; i < cCount; i++ {
			v, _ := strconv.Atoi(colParts[i+1])
			cols = append(cols, v)
		}
	}
	// apply flips
	res := make([][]int, n)
	for i := range res {
		res[i] = make([]int, m)
		copy(res[i], mat[i])
	}
	rowFlip := make([]bool, n)
	colFlip := make([]bool, m)
	for _, r := range rows {
		if r < 1 || r > n {
			return fmt.Errorf("row index out of range")
		}
		if rowFlip[r-1] {
			return fmt.Errorf("duplicate row")
		}
		rowFlip[r-1] = true
	}
	for _, c := range cols {
		if c < 1 || c > m {
			return fmt.Errorf("column index out of range")
		}
		if colFlip[c-1] {
			return fmt.Errorf("duplicate column")
		}
		colFlip[c-1] = true
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := res[i][j]
			if rowFlip[i] {
				val = -val
			}
			if colFlip[j] {
				val = -val
			}
			res[i][j] = val
		}
	}
	// check sums
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			sum += res[i][j]
		}
		if sum < 0 {
			return fmt.Errorf("row %d sum negative", i+1)
		}
	}
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += res[i][j]
		}
		if sum < 0 {
			return fmt.Errorf("column %d sum negative", j+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for {
		n, m, mat, ok := parseMatrix(scanner)
		if !ok {
			break
		}
		idx++
		if err := runCase(bin, n, m, mat); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
