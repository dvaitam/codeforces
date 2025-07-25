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

func expected(n, m, k int, row []int, orders [][]int) string {
	r := make([]int, k)
	copy(r, row)
	total := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := orders[i][j]
			pos := 0
			for pos < k && r[pos] != x {
				pos++
			}
			total += pos + 1
			tmp := r[pos]
			for t := pos; t > 0; t-- {
				r[t] = r[t-1]
			}
			r[0] = tmp
		}
	}
	return fmt.Sprintf("%d", total)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		vals := make([]int, 3)
		for i := 0; i < 3; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			vals[i], _ = strconv.Atoi(scan.Text())
		}
		n, m, k := vals[0], vals[1], vals[2]
		row := make([]int, k)
		for i := 0; i < k; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			row[i], _ = strconv.Atoi(scan.Text())
		}
		orders := make([][]int, n)
		for i := 0; i < n; i++ {
			orders[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if !scan.Scan() {
					fmt.Println("bad test file")
					os.Exit(1)
				}
				orders[i][j], _ = strconv.Atoi(scan.Text())
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < k; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(row[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(orders[i][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := expected(n, m, k, row, orders) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
