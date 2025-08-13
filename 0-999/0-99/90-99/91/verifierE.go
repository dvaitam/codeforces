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
	cmd := exec.Command(bin)
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

func checkCase(lines []string, idx int, bin string) error {
	input := strings.Join(lines, "\n") + "\n"
	got, err := runProg(bin, input)
	if err != nil {
		return fmt.Errorf("case %d failed: %v", idx, err)
	}
	// parse input
	r := bufio.NewReader(strings.NewReader(input))
	var n, q int
	if _, err := fmt.Fscan(r, &n, &q); err != nil {
		return fmt.Errorf("case %d: %v", idx, err)
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i], &b[i])
	}
	L := make([]int, q)
	R := make([]int, q)
	T := make([]int64, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(r, &L[i], &R[i], &T[i])
	}
	// parse candidate output
	outs := strings.Fields(got)
	if len(outs) != q {
		return fmt.Errorf("case %d: expected %d lines of output, got %d", idx, q, len(outs))
	}
	for i := 0; i < q; i++ {
		ans, err := strconv.Atoi(outs[i])
		if err != nil {
			return fmt.Errorf("case %d line %d: %v", idx, i+1, err)
		}
		if ans < L[i] || ans > R[i] {
			return fmt.Errorf("case %d line %d: answer %d out of range [%d,%d]", idx, i+1, ans, L[i], R[i])
		}
		maxH := a[L[i]-1] + b[L[i]-1]*T[i]
		for j := L[i]; j <= R[i]; j++ {
			h := a[j-1] + b[j-1]*T[i]
			if h > maxH {
				maxH = h
			}
		}
		candH := a[ans-1] + b[ans-1]*T[i]
		if candH != maxH {
			return fmt.Errorf("case %d line %d: answer %d yields height %d, expected %d", idx, i+1, ans, candH, maxH)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) > 0 {
				idx++
				if err := checkCase(lines, idx, bin); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				lines = nil
			}
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		idx++
		if err := checkCase(lines, idx, bin); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
