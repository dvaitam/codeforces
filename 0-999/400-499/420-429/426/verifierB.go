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

func mapIndex(i, total, x int) int {
	cur := i
	length := total
	for length > x {
		half := length / 2
		if cur >= half {
			cur = length - cur - 1
		}
		length = half
	}
	return cur
}

func check(a [][]int, total, x int) bool {
	n := total
	for i := x; i < n; i++ {
		idx := mapIndex(i, n, x)
		for j := range a[i] {
			if a[i][j] != a[idx][j] {
				return false
			}
		}
	}
	return true
}

func expected(a [][]int) int {
	n := len(a)
	maxK := 0
	for t := n; t%2 == 0; t /= 2 {
		maxK++
	}
	for k := maxK; k >= 0; k-- {
		div := 1 << k
		if n%div != 0 {
			continue
		}
		x := n / div
		if check(a, n, x) {
			return x
		}
	}
	return n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m int
		fmt.Sscan(line, &n, &m)
		matrix := make([][]int, n)
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "case %d malformed: missing row %d\n", idx, i+1)
				os.Exit(1)
			}
			rowLine := strings.TrimSpace(scanner.Text())
			fields := strings.Fields(rowLine)
			if len(fields) != m {
				fmt.Fprintf(os.Stderr, "case %d malformed: row %d expected %d values got %d\n", idx, i+1, m, len(fields))
				os.Exit(1)
			}
			row := make([]int, m)
			for j, f := range fields {
				v, err := strconv.Atoi(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d malformed: %v\n", idx, err)
					os.Exit(1)
				}
				row[j] = v
			}
			matrix[i] = row
			rows[i] = strings.Join(fields, " ")
		}
		input := line + "\n" + strings.Join(rows, "\n") + "\n"
		want := fmt.Sprintf("%d", expected(matrix))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx, want, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
