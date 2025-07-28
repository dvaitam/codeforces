package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testC struct {
	n int
	u []int
	s []int64
}

func solveC(tc testC) []int64 {
	n := tc.n
	groups := make(map[int][]int64)
	for i := 0; i < n; i++ {
		groups[tc.u[i]] = append(groups[tc.u[i]], tc.s[i])
	}
	prefix := make(map[int][]int64)
	for _, arr := range groups {
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		for i := 1; i < len(arr); i++ {
			arr[i] += arr[i-1]
		}
		size := len(arr)
		ps := prefix[size]
		if ps == nil {
			ps = make([]int64, size)
			prefix[size] = ps
		}
		for i := 0; i < size; i++ {
			ps[i] += arr[i]
		}
	}
	ans := make([]int64, n)
	for size, ps := range prefix {
		for k := 1; k <= size; k++ {
			t := size / k
			if t == 0 {
				break
			}
			idx := t*k - 1
			ans[k-1] += ps[idx]
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	in := bufio.NewReader(file)
	var t int
	fmt.Fscan(in, &t)
	for idx := 1; idx <= t; idx++ {
		var n int
		fmt.Fscan(in, &n)
		u := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &u[i])
		}
		s := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}
		tc := testC{n, u, s}
		expect := solveC(tc)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", u[i])
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", s[i])
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		gotFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(gotFields) != n {
			fmt.Printf("case %d expected %d numbers got %d\n", idx, n, len(gotFields))
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			got, _ := strconv.ParseInt(gotFields[i], 10, 64)
			if got != expect[i] {
				fmt.Printf("case %d failed at k=%d expected %d got %d\n", idx, i+1, expect[i], got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
