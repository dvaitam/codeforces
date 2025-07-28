package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const mod = 998244353

func countChaotic(x, y string) int {
	memo := make(map[[3]int]int)
	var dfs func(i, j int, last byte) int
	dfs = func(i, j int, last byte) int {
		if i == len(x) && j == len(y) {
			return 1
		}
		key := [3]int{i, j, int(last)}
		if v, ok := memo[key]; ok {
			return v
		}
		res := 0
		if i < len(x) && (last == 0 || last != x[i]) {
			res += dfs(i+1, j, x[i])
		}
		if j < len(y) && (last == 0 || last != y[j]) {
			res += dfs(i, j+1, y[j])
		}
		memo[key] = res % mod
		return memo[key]
	}
	return dfs(0, 0, 0)
}

func solve(x, y string) int {
	ans := 0
	for l1 := 0; l1 < len(x); l1++ {
		for r1 := l1 + 1; r1 <= len(x); r1++ {
			xs := x[l1:r1]
			for l2 := 0; l2 < len(y); l2++ {
				for r2 := l2 + 1; r2 <= len(y); r2++ {
					ys := y[l2:r2]
					ans = (ans + countChaotic(xs, ys)) % mod
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		x := scan.Text()
		scan.Scan()
		y := scan.Text()
		expected[i] = fmt.Sprintf("%d", solve(x, y))
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
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
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
