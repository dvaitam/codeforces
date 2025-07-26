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

func expectedC(n, m int, segs [][2]int) int {
	diff := make([]int, m+2)
	for _, s := range segs {
		l, r := s[0], s[1]
		if l < 1 {
			l = 1
		}
		if r > m {
			r = m
		}
		diff[l]++
		if r+1 <= m {
			diff[r+1]--
		}
	}
	cnt := make([]int, m+1)
	cur := 0
	for i := 1; i <= m; i++ {
		cur += diff[i]
		cnt[i] = cur
	}
	lis := make([]int, m+1)
	tails := []int{}
	for i := 1; i <= m; i++ {
		x := cnt[i]
		idx := sort.Search(len(tails), func(j int) bool { return tails[j] > x })
		if idx == len(tails) {
			tails = append(tails, x)
		} else {
			tails[idx] = x
		}
		lis[i] = idx + 1
	}
	lds := make([]int, m+2)
	tails = tails[:0]
	for i := m; i >= 1; i-- {
		x := cnt[i]
		idx := sort.Search(len(tails), func(j int) bool { return tails[j] > x })
		if idx == len(tails) {
			tails = append(tails, x)
		} else {
			tails[idx] = x
		}
		lds[i] = idx + 1
	}
	ans := 0
	for i := 1; i <= m; i++ {
		if v := lis[i] + lds[i] - 1; v > ans {
			ans = v
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+2*n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, 2+2*n, len(parts))
			os.Exit(1)
		}
		segs := make([][2]int, n)
		for i := 0; i < n; i++ {
			l, _ := strconv.Atoi(parts[2+2*i])
			r, _ := strconv.Atoi(parts[3+2*i])
			segs[i] = [2]int{l, r}
		}
		expect := strconv.Itoa(expectedC(n, m, segs))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
