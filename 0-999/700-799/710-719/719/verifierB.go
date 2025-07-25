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

func expected(s string) int {
	n := len(s)
	var cost1_b, cost1_r int // pattern starting with 'r'
	var cost2_b, cost2_r int // pattern starting with 'b'
	for i := 0; i < n; i++ {
		c := s[i]
		if i%2 == 0 {
			if c == 'b' {
				cost1_b++
			} else {
				cost2_r++
			}
		} else {
			if c == 'r' {
				cost1_r++
			} else {
				cost2_b++
			}
		}
	}
	cost1 := cost1_b
	if cost1_r > cost1 {
		cost1 = cost1_r
	}
	cost2 := cost2_b
	if cost2_r > cost2 {
		cost2 = cost2_r
	}
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for idx := 1; idx <= t; idx++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test format at case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "missing string at case %d\n", idx)
			os.Exit(1)
		}
		s := scan.Text()
		if len(s) != n {
			fmt.Fprintf(os.Stderr, "case %d: length mismatch\n", idx)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		exp := expected(s)
		if got != exp {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
