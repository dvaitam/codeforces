package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func solveCase(n, k int, model []string) []string {
	cur := make([][]byte, n)
	for i := 0; i < n; i++ {
		cur[i] = []byte(model[i])
	}
	size := n
	for step := 2; step <= k; step++ {
		nextSize := size * n
		next := make([][]byte, nextSize)
		for i := range next {
			next[i] = make([]byte, nextSize)
			for j := range next[i] {
				next[i][j] = '.'
			}
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if cur[i][j] == '*' {
					for di := 0; di < n; di++ {
						for dj := 0; dj < n; dj++ {
							next[i*n+di][j*n+dj] = '*'
						}
					}
				} else {
					for di := 0; di < n; di++ {
						for dj := 0; dj < n; dj++ {
							next[i*n+di][j*n+dj] = model[di][dj]
						}
					}
				}
			}
		}
		cur = next
		size = nextSize
	}
	res := make([]string, size)
	for i := 0; i < size; i++ {
		res[i] = string(cur[i])
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n, k int
		fmt.Sscan(scan.Text(), &n, &k)
		model := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			model[j] = scan.Text()
		}
		expected := solveCase(n, k, model)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for _, row := range model {
			input.WriteString(row)
			input.WriteByte('\n')
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanLines)
		for _, exp := range expected {
			if !outScan.Scan() {
				fmt.Printf("case %d: missing output line\n", i+1)
				os.Exit(1)
			}
			if outScan.Text() != exp {
				fmt.Printf("case %d mismatch\nexpected: %s\ngot: %s\n", i+1, exp, outScan.Text())
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("case %d: extra output\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
