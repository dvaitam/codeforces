package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

func expected(n, m int, vals []int) int {
	lcnt, rcnt := 0, 0
	posMap := make(map[int]bool)
	for _, v := range vals {
		if v == -1 {
			lcnt++
		} else if v == -2 {
			rcnt++
		} else if v > 0 {
			if v >= 1 && v <= m {
				posMap[v] = true
			}
		}
	}
	pos := make([]int, 0, len(posMap))
	for k := range posMap {
		pos = append(pos, k)
	}
	sort.Ints(pos)
	k := len(pos)
	ans := 0
	if tmp := lcnt + k; tmp > ans {
		if tmp > m {
			tmp = m
		}
		if tmp > ans {
			ans = tmp
		}
	}
	if tmp := rcnt + k; tmp > ans {
		if tmp > m {
			tmp = m
		}
		if tmp > ans {
			ans = tmp
		}
	}
	for i, p := range pos {
		left := p - 1
		if left > lcnt+i {
			left = lcnt + i
		}
		right := m - p
		if right > rcnt+(k-i-1) {
			right = rcnt + (k - i - 1)
		}
		cur := 1 + left + right
		if cur > ans {
			ans = cur
		}
	}
	if ans > m {
		ans = m
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		vals := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			vals[j] = v
		}
		ans := expected(n, m, vals)
		expectedOut[i] = fmt.Sprintf("%d", ans)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expectedOut[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
