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
	const testcasesRaw = `4 7 1 4 2 5 5 6 2 5
5 17 16 16 5 6 1 3 1 8 13 16
3 14 6 12 5 6 13 13
5 8 4 4 1 2 8 8 8 8 5 5
5 5 3 5 3 4 3 5 3 3 1 4
2 20 13 18 16 18
4 5 4 5 3 4 5 5 4 5
5 10 10 10 4 7 4 6 9 9 9 9
5 17 17 17 11 13 14 16 15 17 17 17
2 9 8 8 1 4
5 4 3 3 2 2 2 4 3 4 1 1
1 20 6 19
3 8 7 8 6 7 2 6
4 4 1 2 1 3 3 3 4 4
3 6 2 5 6 6 5 6
5 8 1 1 7 7 3 7 2 5 5 7
5 7 6 6 3 4 3 5 4 7 5 6
3 3 2 2 1 3 2 3
5 15 7 9 10 13 13 13 11 11 10 10
2 16 10 15 4 14
4 9 5 5 6 6 1 1 1 3
1 14 7 14
1 11 11 11
1 13 11 13
2 2 1 2 1 1
5 13 3 11 12 13 7 10 9 9 12 13
5 5 5 5 2 3 1 5 4 5 3 4
2 13 11 11 11 11
5 3 1 3 2 2 1 3 1 1 2 2
3 5 3 5 1 1 3 3
3 15 6 13 8 9 2 7
4 13 13 13 9 12 9 12 9 9
2 5 4 5 2 2
2 5 1 4 5 5
1 17 9 13
3 16 13 15 15 16 8 16
1 13 13 13
3 17 1 8 3 10 4 9
4 20 8 16 1 17 17 17 2 20
3 2 2 2 2 2 2 2
4 1 1 1 1 1 1 1 1 1
3 19 16 17 1 15 6 15
4 11 10 11 7 11 11 11 9 10
4 8 6 7 5 6 4 6 4 8
5 12 8 9 10 11 10 12 3 8 4 6
3 18 6 14 4 18 15 18
2 6 5 6 5 5
1 18 16 16
3 13 5 13 4 9 7 9
1 2 1 1
5 6 5 6 1 3 5 6 6 6 4 5
4 2 1 1 2 2 2 2 2 2
4 12 4 5 12 12 11 11 9 9
5 17 11 16 4 13 8 8 2 5 10 11
5 13 12 12 13 13 5 9 3 5 10 11
2 5 2 4 4 5
2 20 9 12 3 7
4 6 2 6 1 3 5 6 6 6
3 3 3 3 3 3 1 1
2 19 1 14 5 7
4 18 2 18 11 14 6 13 11 18
2 9 4 4 6 9
3 12 9 11 4 8 9 12
1 6 4 4
5 4 2 4 4 4 4 4 2 2 1 3
4 11 7 11 1 6 7 10 9 9
5 19 5 15 2 11 15 19 19 19 7 14
2 13 2 3 9 12
1 17 15 17
2 2 1 1 1 1
5 13 7 7 9 10 9 10 4 11 2 5
1 17 5 14
2 18 17 18 13 18
3 16 1 12 12 13 2 5
4 7 5 7 2 2 2 3 4 7
1 20 9 19
3 13 6 10 11 11 10 10
4 12 10 12 11 11 4 6 7 8
1 16 10 11
5 12 12 12 2 10 6 8 11 11 6 12
2 7 3 7 6 7
4 13 1 3 8 8 5 7 9 10
5 16 12 13 5 14 14 15 7 8 6 8
4 18 14 15 2 11 16 18 3 14
1 4 4 4
5 2 1 1 2 2 2 2 1 1 1 1
4 8 1 4 3 8 3 8 3 7
5 6 1 1 3 3 3 6 5 6 5 6
1 3 1 2
1 12 6 8
4 13 1 5 11 11 4 5 10 11
2 17 7 14 11 17
1 9 7 9
5 18 5 16 11 17 7 11 4 11 9 12
3 10 2 7 3 10 1 7
4 2 2 2 2 2 1 1 2 2
4 9 4 7 8 9 9 9 3 7
4 9 3 8 1 4 6 8 2 2
4 10 6 10 4 6 7 8 2 5
3 16 15 15 12 15 14 15`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
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
