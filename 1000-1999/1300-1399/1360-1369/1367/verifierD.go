package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solveCase(s string, m int, b []int) string {
	freq := make([]int, 26)
	for _, ch := range s {
		freq[int(ch-'a')]++
	}
	res := bytes.Repeat([]byte{'.'}, m)
	used := make([]bool, m)
	remaining := m
	ch := 25
	for remaining > 0 {
		zeros := make([]int, 0)
		for i := 0; i < m; i++ {
			if !used[i] && b[i] == 0 {
				zeros = append(zeros, i)
			}
		}
		if len(zeros) == 0 {
			break
		}
		for ch >= 0 && freq[ch] < len(zeros) {
			ch--
		}
		if ch < 0 {
			break
		}
		for _, pos := range zeros {
			res[pos] = byte('a' + ch)
			used[pos] = true
		}
		freq[ch] -= len(zeros)
		ch--
		remaining -= len(zeros)
		for i := 0; i < m; i++ {
			if used[i] {
				continue
			}
			sum := 0
			for _, pos := range zeros {
				if i > pos {
					sum += i - pos
				} else {
					sum += pos - i
				}
			}
			b[i] -= sum
		}
	}
	return string(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]struct {
		s string
		m int
		b []int
	}, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		s := scan.Text()
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		arr := make([]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		cases[i] = struct {
			s string
			m int
			b []int
		}{s, m, arr}
	}
	expected := make([]string, t)
	for i := range cases {
		expected[i] = solveCase(cases[i].s, cases[i].m, append([]int(nil), cases[i].b...))
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
