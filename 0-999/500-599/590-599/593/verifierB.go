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

type line struct {
	y1 int64
	y2 int64
}

func solve(n int, x1, x2 int64, ks, bs []int64) string {
	lines := make([]line, n)
	for i := 0; i < n; i++ {
		lines[i].y1 = ks[i]*x1 + bs[i]
		lines[i].y2 = ks[i]*x2 + bs[i]
	}
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].y1 == lines[j].y1 {
			return lines[i].y2 < lines[j].y2
		}
		return lines[i].y1 < lines[j].y1
	})
	for i := 0; i < n-1; i++ {
		if lines[i].y2 > lines[i+1].y2 {
			return "Yes"
		}
	}
	return "No"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	type input struct {
		n      int
		x1, x2 int64
		ks, bs []int64
	}
	inputs := make([]input, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x1, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		x2, _ := strconv.ParseInt(scan.Text(), 10, 64)
		ks := make([]int64, n)
		bs := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			ks[j], _ = strconv.ParseInt(scan.Text(), 10, 64)
			scan.Scan()
			bs[j], _ = strconv.ParseInt(scan.Text(), 10, 64)
		}
		inputs[i] = input{n, int64(x1), int64(x2), ks, bs}
	}
	expected := make([]string, t)
	for i, in := range inputs {
		expected[i] = solve(in.n, in.x1, in.x2, in.ks, in.bs)
	}
	for i, in := range inputs {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", in.n, in.x1, in.x2))
		for j := 0; j < in.n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", in.ks[j], in.bs[j]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
