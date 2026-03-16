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

const testcasesBRaw = `100
1 2 7 6 2 3 7 1 5 4 1 2 
3 1 2 1 2 2 2 2 
3 2 2 2 1 2 1 1 2 2 1 
5 1 1 1 1 1 1 1 1 
5 2 2 2 1 2 1 2 1 2 1 1 2 2 1 
3 5 6 4 3 1 2 6 5 2 3 6 1 4 3 1 4 5 2 5 4 3 2 6 
5 1 1 1 1 1 1 1 1 
2 1 4 2 4 1 3 3 4 
4 1 4 1 2 4 3 2 1 4 4 
5 4 4 3 2 4 1 3 2 4 1 1 4 3 2 4 1 3 2 4 1 3 2 4 3 1 2 
1 1 1 1 1 
1 1 3 2 3 1 3 
3 2 8 5 4 3 1 8 6 2 7 8 1 7 1 4 7 
5 2 2 2 1 2 1 1 2 1 2 1 2 2 1 
4 6 6 3 6 1 2 5 4 1 6 4 2 5 3 6 1 2 5 4 3 3 2 1 4 6 5 4 6 5 3 2 1 
1 1 1 1 1 
3 3 4 3 4 1 2 4 1 3 4 2 1 3 4 1 
3 7 8 1 3 5 6 2 8 7 4 2 6 4 5 1 8 3 8 4 7 2 3 5 6 7 2 6 5 8 3 4 
4 3 3 2 3 1 3 2 1 3 1 2 2 1 3 1 3 2 
3 1 4 3 1 2 4 4 1 3 
5 3 3 3 2 1 2 1 3 1 2 3 3 2 1 2 1 3 3 1 2 
4 4 8 7 5 8 3 6 1 2 4 7 8 3 1 3 4 1 7 7 5 1 6 8 5 1 4 
3 1 1 1 1 1 1 
2 1 2 1 2 2 1 
4 1 1 1 1 1 1 1 
2 2 3 2 3 1 2 1 1 2 
4 6 6 6 1 5 2 4 3 4 5 1 3 2 6 1 4 5 2 6 3 2 5 6 3 1 4 4 3 5 2 6 1 
5 2 2 1 2 2 1 2 1 2 1 1 2 1 2 
3 2 2 2 1 2 1 2 1 1 2 
4 4 4 3 1 2 4 2 3 1 4 2 3 4 1 4 2 3 1 1 3 2 4 
3 2 4 2 4 3 1 2 3 4 2 3 4 
2 2 2 1 2 2 1 2 1 
4 1 5 2 4 5 1 3 5 5 4 5 
3 5 7 6 1 7 2 3 5 4 4 5 3 1 6 2 4 7 3 1 4 6 2 3 5 
2 2 4 1 2 4 3 2 1 1 4 
4 2 4 4 2 3 1 4 1 1 3 2 4 1 2 
1 5 8 5 1 4 6 3 2 7 8 3 6 2 7 1 
5 3 4 3 1 2 4 4 3 1 4 2 1 2 1 3 1 2 3 1 2 3 
2 3 6 3 4 6 1 2 5 3 2 1 5 4 6 
4 3 7 3 1 6 5 2 4 7 7 4 1 2 6 1 4 5 1 1 7 6 
1 2 2 2 1 2 1 
4 1 1 1 1 1 1 1 
2 4 7 3 7 1 5 4 6 2 5 6 2 4 3 1 2 6 
3 1 7 5 2 4 3 1 7 6 1 1 1 
4 2 2 2 1 1 2 2 1 1 2 1 2 
3 1 4 3 4 2 1 1 2 3 
5 4 5 3 4 5 2 1 1 5 4 3 1 2 4 3 3 5 4 1 3 1 5 4 3 2 5 4 
5 1 1 1 1 1 1 1 1 
2 1 2 2 1 2 2 
3 1 1 1 1 1 1 
1 1 2 2 1 1 
5 1 1 1 1 1 1 1 1 
4 4 6 3 5 2 4 1 6 4 3 2 5 1 6 3 5 6 5 4 2 5 6 3 4 
5 3 3 3 1 2 1 2 3 2 3 1 2 3 1 3 2 1 1 3 2 
4 2 4 3 1 4 2 2 3 1 2 3 1 3 1 
2 6 6 6 5 1 4 2 3 3 5 2 4 1 6 6 3 1 5 4 2 
5 1 2 2 1 2 1 2 1 1 
2 2 2 2 1 1 2 2 1 
1 1 1 1 1 
2 2 3 1 2 3 1 2 3 2 
5 1 7 3 6 1 5 4 2 7 4 5 3 6 3 
5 1 2 2 1 1 1 2 2 2 
4 1 7 3 5 1 6 2 7 4 6 4 3 1 
2 2 6 6 1 5 4 3 2 6 3 2 3 
1 1 5 1 2 4 5 3 1 
5 4 8 7 2 5 4 6 3 1 8 8 6 4 2 4 8 7 5 2 8 5 3 7 1 6 3 5 8 2 7 
3 6 7 3 2 5 6 4 1 7 5 6 1 4 7 2 7 6 1 5 3 4 2 5 7 4 6 1 
4 7 8 5 4 6 7 8 3 2 1 8 3 5 4 2 1 7 1 3 4 8 7 5 6 7 8 6 3 5 2 1 2 3 1 8 4 5 7 
4 2 2 2 1 2 1 2 1 1 2 2 1 
2 4 5 1 2 5 4 3 3 1 2 5 4 5 3 2 
2 1 1 1 1 1 
3 2 5 3 2 5 4 1 1 5 5 1 3 5 
2 2 3 1 2 3 1 3 3 1 
4 6 7 7 1 5 6 2 3 4 3 4 7 6 2 5 6 3 2 1 7 5 6 1 3 7 4 5 5 6 4 7 3 1 
4 3 5 4 5 2 1 3 3 4 2 2 1 4 4 5 2 1 5 2 
4 5 6 2 3 4 6 1 5 3 5 6 2 4 4 1 3 2 5 1 6 4 5 2 4 5 3 6 2 
1 1 8 4 1 7 2 6 3 8 5 7 
4 5 7 7 6 3 4 2 1 5 5 6 4 2 3 1 3 5 2 6 6 4 3 5 2 6 2 4 1 5 
4 5 5 1 5 2 4 3 4 3 5 2 1 2 1 4 5 3 4 2 1 5 3 4 1 5 2 3 
3 2 2 2 1 2 1 1 2 2 1 
1 5 5 1 4 2 5 3 3 4 1 5 2 
2 4 7 2 3 5 1 6 4 7 2 6 3 7 6 2 1 7 
4 1 1 1 1 1 1 1 
2 2 3 2 1 3 2 1 2 3 
4 2 3 3 2 1 1 2 3 1 3 2 2 1 
3 6 7 6 5 3 7 1 4 2 4 7 1 2 6 3 3 1 4 7 6 2 3 6 1 5 7 4 
2 1 4 2 1 4 3 4 4 
5 1 4 3 1 2 4 2 4 1 4 4 
1 1 7 7 4 5 3 2 1 6 2 
1 1 1 1 1 
5 2 7 7 4 2 6 3 1 5 6 3 2 1 7 5 5 3 3 1 
2 6 7 4 7 1 2 5 6 3 4 5 3 1 7 2 3 6 5 4 2 7 
4 3 7 2 5 3 1 7 6 4 7 4 6 4 2 1 1 6 3 6 4 7 
2 1 5 4 3 1 2 5 2 5 
5 3 4 4 2 3 1 3 4 2 3 4 2 4 1 3 3 1 4 2 3 1 
2 2 4 1 2 3 4 1 3 2 4 
3 4 6 5 4 3 2 1 6 3 5 4 6 2 1 5 6 1 4 6 2 
1 2 3 1 3 2 1 2 
2 1 2 2 1 1 1 
5 4 4 3 4 2 1 1 3 2 4 4 2 1 3 3 4 2 1 3 2 4 1 2 3 1 4 `

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data := []byte(testcasesBRaw)
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
