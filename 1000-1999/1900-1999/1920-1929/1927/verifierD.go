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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		if len(fields) < pos+n {
			fmt.Printf("test %d short array\n", idx)
			os.Exit(1)
		}
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			v, _ := strconv.Atoi(fields[pos+i])
			a[i] = v
			sb.WriteString(fields[pos+i])
		}
		sb.WriteByte('\n')
		pos += n
		if len(fields) <= pos {
			fmt.Printf("test %d missing q\n", idx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(fields[pos])
		pos++
		sb.WriteString(fmt.Sprintf("%d\n", q))
		if len(fields) != pos+2*q {
			fmt.Printf("test %d wrong query count\n", idx)
			os.Exit(1)
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, _ := strconv.Atoi(fields[pos+2*i])
			r, _ := strconv.Atoi(fields[pos+2*i+1])
			queries[i] = [2]int{l, r}
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		input := sb.String()
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		outScanner := bufio.NewScanner(strings.NewReader(got))
		answers := make([]string, 0, q)
		for outScanner.Scan() {
			line := strings.TrimSpace(outScanner.Text())
			if line != "" {
				answers = append(answers, line)
			}
		}
		if err := outScanner.Err(); err != nil {
			fmt.Printf("test %d output scan error: %v\n", idx, err)
			os.Exit(1)
		}
		if len(answers) != q {
			fmt.Printf("test %d expected %d lines, got %d\n", idx, q, len(answers))
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			parts := strings.Fields(answers[i])
			if len(parts) != 2 {
				fmt.Printf("test %d line %d invalid output\n", idx, i+1)
				os.Exit(1)
			}
			x, err1 := strconv.Atoi(parts[0])
			y, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				fmt.Printf("test %d line %d non-integer output\n", idx, i+1)
				os.Exit(1)
			}
			l, r := queries[i][0], queries[i][1]
			if x == -1 && y == -1 {
				same := true
				for k := l; k < r; k++ {
					if a[k-1] != a[k] {
						same = false
						break
					}
				}
				if !same {
					fmt.Printf("test %d line %d incorrect -1 -1\n", idx, i+1)
					os.Exit(1)
				}
			} else {
				if x < l || x > r || y < l || y > r {
					fmt.Printf("test %d line %d indices out of range\n", idx, i+1)
					os.Exit(1)
				}
				if a[x-1] == a[y-1] {
					fmt.Printf("test %d line %d values not distinct\n", idx, i+1)
					os.Exit(1)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
