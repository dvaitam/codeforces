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

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(arr []int) int {
	n := len(arr)
	sorted := true
	for i := 0; i < n; i++ {
		if arr[i] != i+1 {
			sorted = false
			break
		}
	}
	if sorted {
		return 0
	}
	l := 0
	for l < n && arr[l] == l+1 {
		l++
	}
	r := n - 1
	for r >= 0 && arr[r] == r+1 {
		r--
	}
	needTwo := false
	for i := l; i <= r; i++ {
		if arr[i] == i+1 {
			needTwo = true
			break
		}
	}
	if needTwo {
		return 2
	}
	return 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][]int, t)
	for idx := 0; idx < t; idx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = v
		}
		cases[idx] = arr
	}
	out, err := runCandidate(bin, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(strings.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, arr := range cases {
		expect := solveCase(arr)
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(outScan.Text())
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
