package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func bestScore(a, b int) (int64, string) {
	if a == 0 {
		score := -int64(b) * int64(b)
		return score, strings.Repeat("x", b)
	}
	if b == 0 {
		score := int64(a) * int64(a)
		return score, strings.Repeat("o", a)
	}
	bestS := int64(-1 << 63)
	bestK := 1
	for k := 1; k <= a; k++ {
		big := int64(a - k + 1)
		sumO := big*big + int64(k-1)
		m := k + 1
		if b < m {
			m = b
		}
		var sumX int64
		if m > 0 {
			base := b / m
			r := b % m
			sumX += int64(r) * int64(base+1) * int64(base+1)
			sumX += int64(m-r) * int64(base) * int64(base)
		}
		s := sumO - sumX
		if s > bestS {
			bestS = s
			bestK = k
		}
	}
	k := bestK
	oLens := make([]int, k)
	oLens[0] = a - k + 1
	for i := 1; i < k; i++ {
		oLens[i] = 1
	}
	slots := make([]int, k+1)
	m := k + 1
	if b < m {
		m = b
	}
	if m > 0 {
		base := b / m
		r := b % m
		for i := 0; i < m; i++ {
			if i < r {
				slots[i] = base + 1
			} else {
				slots[i] = base
			}
		}
	}
	var sb strings.Builder
	for i := 0; i <= k; i++ {
		if slots[i] > 0 {
			sb.WriteString(strings.Repeat("x", slots[i]))
		}
		if i < k {
			sb.WriteString(strings.Repeat("o", oLens[i]))
		}
	}
	return bestS, sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
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
		var a, b int
		fmt.Sscan(line, &a, &b)
		expScore, expStr := bestScore(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outs := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(outs) < 2 {
			fmt.Printf("Test %d invalid output: %s\n", idx, string(out))
			os.Exit(1)
		}
		gotScore, err := strconv.ParseInt(strings.TrimSpace(outs[0]), 10, 64)
		if err != nil {
			fmt.Printf("Test %d: first line not integer\n", idx)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(outs[1])
		if gotScore != expScore || gotStr != expStr {
			fmt.Printf("Test %d failed\nexpected:\n%d\n%s\n\ngot:\n%s", idx, expScore, expStr, string(out))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
