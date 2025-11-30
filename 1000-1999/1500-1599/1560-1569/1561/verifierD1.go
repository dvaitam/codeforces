package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesD1.txt.
const testcaseData = `
2 998244353
5 100000007
1959 998244353
1943 998244353
117 100000007
175 100000039
1713 100000037
1509 998244353
1373 998244353
633 100000039
250 100000007
828 100000007
1320 998244353
430 998244353
1713 1000000007
715 100000037
1330 100000007
902 100000007
1790 998244353
171 100000039
1418 100000037
1792 100000007
1086 1000000007
54 998244353
1402 100000039
919 998244353
782 1000000007
171 1000000007
1422 100000037
707 100000007
1934 100000039
314 100000007
1125 100000007
1529 100000037
1355 100000007
363 100000039
155 100000037
1920 100000037
1242 100000037
300 998244353
70 100000007
1153 998244353
289 100000039
362 100000039
329 998244353
1274 1000000007
1841 100000007
769 100000037
1276 998244353
638 1000000007
1695 998244353
750 998244353
1005 998244353
911 100000037
1846 1000000007
1292 100000037
1283 100000037
625 100000007
1522 100000037
1455 100000007
92 1000000007
800 1000000007
1526 1000000007
620 998244353
588 100000039
1938 1000000007
344 1000000007
1998 1000000007
1658 100000037
306 100000037
1202 998244353
1648 1000000007
1895 100000037
75 1000000007
661 100000037
556 100000007
1309 100000037
618 100000037
1452 100000039
1799 100000039
913 1000000007
1265 100000037
946 100000037
696 998244353
864 1000000007
1808 100000037
1985 100000007
1691 100000037
348 100000039
375 100000039
1194 998244353
298 998244353
1503 998244353
346 100000039
676 998244353
777 100000039
141 100000007
145 100000039
1230 100000037
1278 1000000007
443 100000039
1976 100000039
760 100000039
1321 100000037
675 998244353
1862 100000037
`

func expectedD1(n, mod int) int {
	dp := make([]int, n+2)
	pref := make([]int, n+2)
	diff := make([]int, n+2)

	dp[1] = 1
	pref[1] = 1
	for j := 2; j <= n; j++ {
		start := j
		end := j + j
		if end > n+1 {
			end = n + 1
		}
		diff[start] = (diff[start] + dp[1]) % mod
		if end <= n {
			diff[end] = (diff[end] - dp[1]) % mod
		}
	}

	add := 0
	for i := 2; i <= n; i++ {
		add = (add + diff[i]) % mod
		if add < 0 {
			add += mod
		}
		dp[i] = (pref[i-1] + add) % mod
		pref[i] = (pref[i-1] + dp[i]) % mod

		for j := 2; i*j <= n; j++ {
			start := i * j
			end := start + j
			if end > n+1 {
				end = n + 1
			}
			diff[start] = (diff[start] + dp[i]) % mod
			if end <= n {
				diff[end] = (diff[end] - dp[i]) % mod
			}
		}
	}

	ans := dp[n] % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func parseTestcases() ([][2]int, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("expected pairs of values, got %d tokens", len(fields))
	}
	res := make([][2]int, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad n at pair %d: %v", i/2+1, err)
		}
		mod, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad mod at pair %d: %v", i/2+1, err)
		}
		res = append(res, [2]int{n, mod})
	}
	return res, nil
}

func runCandidate(bin string, n, mod int) (string, error) {
	input := fmt.Sprintf("%d %d\n", n, mod)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := expectedD1(tc[0], tc[1])
		got, err := runCandidate(bin, tc[0], tc[1])
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strconv.Itoa(expect) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
