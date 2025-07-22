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

func longestGood(arr []int64) int {
	n := len(arr)
	if n <= 2 {
		return n
	}
	cur := 2
	ans := 2
	for i := 2; i < n; i++ {
		if arr[i] == arr[i-1]+arr[i-2] {
			cur++
		} else {
			if cur > ans {
				ans = cur
			}
			cur = 2
		}
	}
	if cur > ans {
		ans = cur
	}
	return ans
}

func runCase(bin string, arr []int64) error {
	var in bytes.Buffer
	fmt.Fprintln(&in, len(arr))
	for i, v := range arr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprint(&in, v)
	}
	in.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", longestGood(arr))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "missing n on case %d\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "missing value on case %d pos %d\n", i+1, j+1)
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = v
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
