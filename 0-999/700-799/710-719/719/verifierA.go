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

func expected(seq []int) string {
	n := len(seq)
	last := seq[n-1]
	if n == 1 {
		switch last {
		case 0:
			return "UP"
		case 15:
			return "DOWN"
		default:
			return "-1"
		}
	}
	if last == 0 {
		return "UP"
	}
	if last == 15 {
		return "DOWN"
	}
	prev := seq[n-2]
	if last > prev {
		return "UP"
	}
	return "DOWN"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
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
		seq := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "incomplete case %d\n", idx)
				os.Exit(1)
			}
			seq[i], _ = strconv.Atoi(scan.Text())
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range seq {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
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
		got := strings.TrimSpace(out.String())
		exp := expected(seq)
		if got != exp {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
