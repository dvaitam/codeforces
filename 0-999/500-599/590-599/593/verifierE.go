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

const testcasesERaw = `100
1 2 2
3 1 2 4
2 1 2 5
3 1 1
3 1 1 5
1 3 2
3 1 3 3
1 1 1 6
2 3 3
3 1 3 5
3 2 1 6
1 1 1 8
3 2 3
2 2 1 4
1 1 2 6
1 3 1 8
2 2 2
2 2 2 5
1 2 2 7
2 3 2
3 1 2 5
3 1 3 7
3 1 4
3 3 1 3
1 3 1 4
1 1 1 6
1 3 1 8
1 2 4
3 1 1 5
1 1 2 7
2 1 1 10
1 1 2 13
1 1 3
3 1 1 5
2 1 1 7
2 1 1 8
1 3 3
2 1 1 3
2 1 3 6
1 1 3 7
3 3 2
3 2 2 3
2 2 2 4
2 3 4
3 2 1 3
3 2 2 6
3 2 2 7
1 1 3 8
2 3 2
1 2 1 5
3 2 1 6
1 1 2
3 1 1 4
2 1 1 7
2 2 1
2 2 1 5
2 3 3
3 2 2 3
3 1 3 4
1 2 1 5
1 3 1
3 1 2 4
2 2 4
1 1 2 5
2 2 2 6
1 1 1 8
3 1 1 10
3 2 5
3 1 2 5
1 1 1 8
2 2 1 9
2 3 2 11
3 2 2 14
2 2 2
3 1 2 3
2 1 1 4
3 1 4
2 3 1 4
3 2 1 7
3 2 1 10
3 1 1 13
2 3 3
3 1 1 3
2 2 3 4
1 1 3 7
2 3 5
1 1 3 3
3 1 1 6
1 1 3 7
1 1 1 8
1 1 2 10
3 2 1
2 3 1 4
3 2 1
3 2 1 3
1 2 3
3 1 1 3
2 1 2 4
2 1 1 7
3 1 1
1 2 1 4
3 3 4
3 3 3 3
1 3 1 4
2 1 1 6
3 1 1 9
3 2 2
2 2 2 3
3 2 1 6
2 3 3
2 2 2 5
2 1 1 7
3 1 1 10
1 1 3
1 1 1 5
1 1 1 7
2 1 1 8
1 3 1
3 1 2 5
1 2 2
1 1 2 5
2 1 1 7
1 1 1
2 1 1 5
1 3 1
1 1 1 3
2 3 3
1 2 1 3
3 1 1 4
1 2 2 6
3 3 2
3 3 2 4
1 1 1 7
1 2 2
1 1 2 3
1 1 2 5
2 1 3
3 2 1 4
2 1 1 5
2 2 1 6
3 2 5
1 3 1 5
3 1 1 8
2 2 1 11
1 2 2 12
1 3 2 14
2 2 4
3 2 2 5
1 1 1 6
1 1 2 9
2 2 2 10
2 1 3
3 2 1 5
3 2 1 8
3 1 1 10
1 3 5
3 1 2 3
3 1 2 5
2 1 3 6
2 1 3 8
1 1 3 11
3 2 1
2 3 2 5
2 2 5
3 2 1 4
1 1 1 7
3 2 1 9
2 1 2 11
3 2 2 12
2 2 1
2 2 2 5
2 1 4
3 2 1 5
1 1 1 7
3 1 1 8
2 1 1 9
2 2 1
3 2 1 4
1 2 2
1 1 1 3
2 1 2 5
2 3 3
3 1 1 3
1 2 2 4
3 2 3 5
3 2 4
1 2 2 4
2 1 2 5
1 2 2 7
2 2 1 10
1 3 2
1 1 1 5
2 1 1 6
1 3 3
2 1 3 5
2 1 1 8
1 1 1 9
2 2 1
2 2 1 3
2 3 4
1 2 1 5
2 2 1 8
2 2 3 10
2 2 2 12
3 2 1
1 1 1 3
3 2 4
1 3 2 5
1 3 1 6
3 2 2 7
3 2 1 8
2 2 2
1 2 1 4
1 1 2 5
3 1 3
2 2 1 5
3 1 1 7
3 3 1 9
3 2 5
1 2 2 5
3 2 2 7
1 2 1 9
2 1 1 10
1 1 2 12
1 2 1
1 1 2 4
2 3 2
2 2 2 4
2 2 3 7
2 3 4
2 1 3 5
1 2 2 6
1 1 1 8
2 2 1 10
1 1 2
1 1 1 4
3 1 1 5
1 2 2
1 1 1 3
1 1 1 4
3 3 1
2 1 2 3
3 2 4
1 1 1 4
2 1 1 7
3 3 2 10
3 2 1 12
3 3 5
3 1 1 5
1 2 3 6
3 2 2 9
3 2 3 11
2 1 2 13
3 2 5
3 1 2 3
1 3 2 6
1 2 2 8
2 3 2 11
3 1 1 12
3 2 2
2 3 1 5
1 1 1 6
2 2 3
1 1 1 5
1 1 1 7
1 2 1 8
1 2 5
1 1 2 3
2 1 1 4
1 1 2 7
3 1 2 8
2 1 1 10
1 3 1
3 1 3 5
2 1 2
1 1 1 5
3 1 1 7
1 1 4
1 1 1 4
3 1 1 6
1 1 1 9
2 1 1 10
3 1 4
3 3 1 4
2 1 1 5
2 3 1 8
3 3 1 10
1 2 2
3 1 1 5
1 1 2 7
1 3 2
1 1 3 5
1 1 3 6
2 1 1
1 2 1 4
1 2 5
3 1 2 4
1 1 2 5
2 1 1 8
3 1 2 10
3 1 1 11
3 3 1
1 2 2 3
1 2 2
2 1 2 5
1 1 2 6
2 2 1
2 1 1 3
1 2 5
3 1 1 4
3 1 2 5
1 1 1 6
2 1 1 7
2 1 2 8
3 3 5
3 3 2 5
3 3 3 7
3 1 3 9
2 1 2 10
3 1 2 12
1 2 5
2 1 1 5
1 1 2 8
2 1 1 9
2 1 1 12
1 1 1 14
2 2 3
3 1 1 3
2 1 1 4
2 1 1 6
2 2 5
1 2 2 4
2 2 2 6
2 1 2 9
3 2 1 10
1 2 2 11
3 3 5
2 2 2 3
1 1 3 5
1 2 2 8
3 2 1 11
2 1 3 13
3 2 5
3 2 1 5
3 3 1 7
1 3 1 9
3 3 1 12
2 2 1 14
2 3 4
2 2 3 4
2 2 1 7
3 1 2 8
2 2 1 11
2 3 1
2 1 1 3
2 2 5
1 2 1 3
2 2 1 6
1 2 1 8
1 1 2 11
3 2 2 12
3 3 1
3 2 1 4
3 3 4
2 1 2 3
1 1 2 6
1 1 3 8
1 2 3 9
2 1 1
1 1 1 3
1 3 3
1 1 3 5
2 1 3 6
2 1 1 9
3 3 4
1 1 3 4
1 3 2 7
1 1 1 10
3 3 2 13
3 3 5
2 3 1 4
3 1 3 6
3 3 1 8
1 1 2 10
2 3 3 12
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesERaw)
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
