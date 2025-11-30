package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `15
01010
1
2 2
0000
3
2 2
13 4
5 5
0
5
29 5
31 3
19 1
17 2
5 3
0
2
17 5
11 2
010
4
5 3
17 2
2 4
3 1
000
3
23 1
19 5
3 4
000
3
31 3
17 2
7 5
00
2
3 3
5 3
0
3
7 3
5 1
2 3
0000
1
2 5
0
3
29 4
11 4
13 5
0000
5
5 5
2 2
23 4
11 2
3 1
000
5
13 2
23 5
5 1
3 5
31 3
000
5
3 3
11 2
29 2
19 4
23 1
011
3
17 3
5 2
29 2
010
1
31 2`

type PA struct {
	p int
	a int
}

func solveCase(s string, arr []PA) string {
	m := len(s)
	zeroPos := make([]int, 0)
	idxMap := make([]int, m)
	for i := 0; i < m; i++ {
		if s[i] == '0' {
			idxMap[i] = len(zeroPos)
			zeroPos = append(zeroPos, i)
		} else {
			idxMap[i] = -1
		}
	}
	z := len(zeroPos)
	fullMask := uint64(1<<z) - 1
	dp := map[uint64]*big.Int{0: big.NewInt(1)}
	for _, pa := range arr {
		p := pa.p
		a := pa.a
		pow := new(big.Int).Exp(big.NewInt(int64(p)), big.NewInt(int64(a-1)), nil)
		patterns := make(map[uint64]int64)
		if p > m {
			patterns[0] = int64(p - m)
			for i := 0; i < m; i++ {
				if s[i] == '0' {
					mask := uint64(1) << uint(idxMap[i])
					patterns[mask]++
				}
			}
		} else {
			used := make([]bool, p)
			masks := make([]uint64, p)
			valid := make([]bool, p)
			for i := range valid {
				valid[i] = true
			}
			for i := 0; i < m; i++ {
				r := (p - (i % p)) % p
				used[r] = true
				if s[i] == '1' {
					valid[r] = false
				} else {
					masks[r] |= 1 << uint(idxMap[i])
				}
			}
			for r := 0; r < p; r++ {
				if !used[r] {
					patterns[0]++
				} else if valid[r] {
					patterns[masks[r]]++
				}
			}
		}
		newdp := make(map[uint64]*big.Int)
		for mask, val := range dp {
			for pat, cnt := range patterns {
				if cnt == 0 {
					continue
				}
				nm := mask | pat
				add := new(big.Int).Mul(val, pow)
				add.Mul(add, big.NewInt(cnt))
				if ex, ok := newdp[nm]; ok {
					ex.Add(ex, add)
				} else {
					newdp[nm] = add
				}
			}
		}
		dp = newdp
	}
	ans := dp[fullMask]
	if ans == nil {
		ans = new(big.Int)
	}
	return ans.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(embeddedTestcases)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := strings.TrimSpace(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		arr := make([]PA, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &arr[i].p, &arr[i].a)
		}
		expected := solveCase(s, arr)
		var input bytes.Buffer
		input.WriteString(s)
		input.WriteByte('\n')
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d %d\n", arr[i].p, arr[i].a)
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
