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

type Pair struct{ x, y int }

func expectedC(n, m int, arr []int) (int, int, []int) {
	a := make([]int, n)
	for i := range arr {
		a[i] = arr[i] - 1
	}
	cnt := make([]int, m)
	for _, v := range a {
		if v >= 0 && v < m {
			cnt[v]++
		}
	}
	per := n / m
	more := n % m
	id := make([]int, m)
	for i := range id {
		id[i] = i
	}
	sort.Slice(id, func(i, j int) bool { return cnt[id[i]] > cnt[id[j]] })
	tar := make([]int, m)
	for idx, v := range id {
		if idx < more {
			tar[v] = per + 1
		} else {
			tar[v] = per
		}
	}
	ans := 0
	for i := 0; i < m; i++ {
		for cnt[i] > tar[i] {
			idx := -1
			for j := 0; j < n; j++ {
				if a[j] == i {
					idx = j
					break
				}
			}
			for j := 0; j < m; j++ {
				if cnt[j] < tar[j] {
					a[idx] = j
					cnt[j]++
					cnt[i]--
					ans++
					break
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		if a[i] < 0 || a[i] >= m {
			for j := 0; j < m; j++ {
				if cnt[j] < tar[j] {
					a[i] = j
					cnt[j]++
					ans++
					break
				}
			}
		}
	}
	out := make([]int, n)
	for i := range a {
		out[i] = a[i] + 1
	}
	return per, ans, out
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 2 {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		nums := strings.Fields(scan.Text())
		if len(nums) != n {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		arr := make([]int, n)
		for j, s := range nums {
			arr[j], _ = strconv.Atoi(s)
		}
		per, ans, outArr := expectedC(n, m, arr)
		input := fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(nums, " "))
		outStrs := make([]string, n)
		for j, v := range outArr {
			outStrs[j] = strconv.Itoa(v)
		}
		exp := fmt.Sprintf("%d %d\n%s\n", per, ans, strings.Join(outStrs, " "))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
