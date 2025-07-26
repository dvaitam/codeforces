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

func expected(n, m int, k int64, arr []int64) int64 {
	var best int64
	for l := 0; l < n; l++ {
		var sum int64
		for r := l; r < n; r++ {
			sum += arr[r]
			length := r - l + 1
			cost := sum - k*int64((length+m-1)/m)
			if cost > best {
				best = cost
			}
		}
	}
	if best < 0 {
		return 0
	}
	return best
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to open testcasesD.txt:", err)
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
		scan.Scan()
		mval, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kval, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[i] = val
		}
		exp := expected(n, mval, int64(kval), append([]int64(nil), arr...))
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", n, mval, kval)
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
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("case %d: invalid output %q\n", caseNum, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
