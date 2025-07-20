package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func expectedB1(n int, m int64, arr []int64) int64 {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	l, r := 0, 1
	sum := arr[0]
	var ans int64
	if arr[0] > m {
		return 0
	}
	for l < n {
		if r < n && arr[r]-arr[l] <= 1 && sum+arr[r] <= m {
			sum += arr[r]
			r++
		} else {
			sum -= arr[l]
			l++
		}
		if sum > ans {
			ans = sum
		}
	}
	if arr[0] > ans {
		ans = arr[0]
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB1.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int
		var m int64
		fmt.Sscan(fields[0], &n)
		fmt.Sscan(fields[1], &m)
		if len(fields) != 2+n {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[2+i], &arr[i])
		}
		expect := expectedB1(n, m, append([]int64(nil), arr...))
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d ", arr[i])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int64
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
