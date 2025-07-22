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

type pair struct{ val, idx int }

func solve(nums []int) string {
	n := len(nums)
	a := make([]pair, n)
	for i, v := range nums {
		a[i] = pair{v, i}
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	ans := make([]int, n)
	q := 0
	for _, p := range a {
		ans[p.idx] = q
		q ^= 1
	}
	var b strings.Builder
	for i, v := range ans {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("invalid testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1]) // m is read but not used in solution
		if len(fields) != 2+n {
			fmt.Printf("invalid testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i], _ = strconv.Atoi(fields[2+i])
		}
		expect := solve(nums)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, m)
		for i, v := range nums {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", idx+1, buf.String(), expect, got)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
