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

type node struct {
	val int64
	pos int
}

func solveCase(n, m, k int, arr []int64) (string, string) {
	b := make([]node, n)
	for i := 0; i < n; i++ {
		b[i] = node{arr[i], i}
	}
	c := make([]node, n)
	copy(c, b)
	sort.Slice(c, func(i, j int) bool { return c[i].val > c[j].val })
	tar := m * k
	vis := make([]bool, n)
	sum := int64(0)
	for i := 0; i < tar; i++ {
		sum += c[i].val
		vis[c[i].pos] = true
	}
	cuts := make([]int, 0, k-1)
	cnt := 0
	t := 0
	for i := 0; i < n && t < k-1; i++ {
		if vis[i] {
			cnt++
		}
		if cnt == m {
			cuts = append(cuts, i+1)
			cnt = 0
			t++
		}
	}
	cutStr := ""
	for i, v := range cuts {
		if i > 0 {
			cutStr += " "
		}
		cutStr += fmt.Sprint(v)
	}
	return fmt.Sprint(sum), cutStr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
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
		if len(fields) < 3 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n := 0
		m := 0
		k := 0
		fmt.Sscan(fields[0], &n)
		fmt.Sscan(fields[1], &m)
		fmt.Sscan(fields[2], &k)
		if len(fields) != 3+n {
			fmt.Printf("test %d expecting %d nums got %d\n", idx, 3+n, len(fields))
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[3+i], &arr[i])
		}
		exp1, exp2 := solveCase(n, m, k, arr)
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		output := strings.Fields(strings.TrimSpace(out.String()))
		if len(output) < 1 {
			fmt.Printf("test %d no output\n", idx)
			os.Exit(1)
		}
		if output[0] != exp1 {
			fmt.Printf("test %d first line expected %s got %s\n", idx, exp1, output[0])
			os.Exit(1)
		}
		if k > 1 {
			gotCut := strings.Join(output[1:], " ")
			if gotCut != exp2 {
				fmt.Printf("test %d second line expected %s got %s\n", idx, exp2, gotCut)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
