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

func expectedQuery(a []int, l, r int) int {
	freq := make(map[int]int)
	for i := l - 1; i < r; i++ {
		freq[a[i]]++
	}
	ans := 0
	for v, c := range freq {
		if c%2 == 0 {
			ans ^= v
		}
	}
	return ans
}

func runCase(exe string, n int, arr []int, queries [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	input := sb.String()
	expVals := make([]int, len(queries))
	for i, q := range queries {
		expVals[i] = expectedQuery(arr, q[0], q[1])
	}
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(expVals) {
		return fmt.Errorf("expected %d numbers got %d", len(expVals), len(fields))
	}
	for i, exp := range expVals {
		got, err := strconv.Atoi(fields[i])
		if err != nil || got != exp {
			return fmt.Errorf("at %d expected %d got %s", i+1, exp, fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for idx := 0; idx < t; idx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			sc.Scan()
			arr[i], _ = strconv.Atoi(sc.Text())
		}
		sc.Scan()
		m, _ := strconv.Atoi(sc.Text())
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			sc.Scan()
			l, _ := strconv.Atoi(sc.Text())
			sc.Scan()
			r, _ := strconv.Atoi(sc.Text())
			queries[i] = [2]int{l, r}
		}
		if err := runCase(exe, n, arr, queries); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
