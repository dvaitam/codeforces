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

const MOD int64 = 1_000_000_007

type event struct {
	tp   int
	x, y int
	t    int
}

type inputE struct {
	n, m   int
	events []event
}

func parseInput(data []byte) ([]inputE, error) {
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	T, _ := strconv.Atoi(scan.Text())
	tests := make([]inputE, T)
	for i := 0; i < T; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		q, _ := strconv.Atoi(scan.Text())
		evs := make([]event, q)
		for j := 0; j < q; j++ {
			scan.Scan()
			tp, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			t, _ := strconv.Atoi(scan.Text())
			evs[j] = event{tp, x, y, t}
		}
		tests[i] = inputE{n, m, evs}
	}
	return tests, nil
}

func step(dp [][]int64, blocked [][]bool) [][]int64 {
	n := len(dp)
	m := len(dp[0])
	nxt := make([][]int64, n)
	for i := 0; i < n; i++ {
		nxt[i] = make([]int64, m)
	}
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			if blocked[r][c] {
				continue
			}
			val := dp[r][c]
			if val == 0 {
				continue
			}
			// stay
			nxt[r][c] = (nxt[r][c] + val) % MOD
			if r > 0 && !blocked[r-1][c] {
				nxt[r-1][c] = (nxt[r-1][c] + val) % MOD
			}
			if r+1 < n && !blocked[r+1][c] {
				nxt[r+1][c] = (nxt[r+1][c] + val) % MOD
			}
			if c > 0 && !blocked[r][c-1] {
				nxt[r][c-1] = (nxt[r][c-1] + val) % MOD
			}
			if c+1 < m && !blocked[r][c+1] {
				nxt[r][c+1] = (nxt[r][c+1] + val) % MOD
			}
		}
	}
	return nxt
}

func solve(inp inputE) []int64 {
	n, m := inp.n, inp.m
	blocked := make([][]bool, n)
	for i := 0; i < n; i++ {
		blocked[i] = make([]bool, m)
	}
	dp := make([][]int64, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int64, m)
	}
	dp[0][0] = 1
	curTime := 1
	var ans []int64
	for _, ev := range inp.events {
		for curTime < ev.t {
			dp = step(dp, blocked)
			curTime++
		}
		if ev.tp == 1 {
			ans = append(ans, dp[ev.x-1][ev.y-1]%MOD)
		} else if ev.tp == 2 {
			blocked[ev.x-1][ev.y-1] = true
		} else {
			blocked[ev.x-1][ev.y-1] = false
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	tests, err := parseInput(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for ti, tst := range tests {
		expected := solve(tst)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tst.n, tst.m, len(tst.events)))
		for _, ev := range tst.events {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", ev.tp, ev.x, ev.y, ev.t))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", ti+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		for i, exp := range expected {
			if !outScan.Scan() {
				fmt.Printf("missing output after query %d in test %d\n", i+1, ti+1)
				os.Exit(1)
			}
			got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
			if got%MOD != exp%MOD {
				fmt.Printf("test %d query %d expected %d got %d\n", ti+1, i+1, exp%MOD, got%MOD)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
