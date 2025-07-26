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

func expected(a []int) string {
	maxVal := a[0]
	idx := 0
	for i := 1; i < len(a); i++ {
		if a[i] > maxVal {
			maxVal = a[i]
			idx = i
		}
	}
	if idx == 0 || idx == len(a)-1 {
		return "NO"
	}
	return "YES"
}

func runExe(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
		return
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad test file")
		return
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			return
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		exp := expected(arr)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		out, err := runExe(bin, input.Bytes())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(out)) != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
