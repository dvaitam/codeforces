package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	mod int64
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
100
70 998244353
293 998244353
435 998244353
412 998244353
393 998244353
34 998244353
132 998244353
62 998244353
255 998244353
391 998244353
232 998244353
243 998244353
335 998244353
196 998244353
405 998244353
109 998244353
50 998244353
251 998244353
16 998244353
459 998244353
429 998244353
201 998244353
223 998244353
313 998244353
392 998244353
394 998244353
3 998244353
358 998244353
230 998244353
138 998244353
371 998244353
412 998244353
119 998244353
304 998244353
485 998244353
54 998244353
463 998244353
164 998244353
17 998244353
13 998244353
15 998244353
334 998244353
279 998244353
6 998244353
482 998244353
453 998244353
197 998244353
353 998244353
112 998244353
498 998244353
218 998244353
373 998244353
16 998244353
272 998244353
115 998244353
393 998244353
226 998244353
482 998244353
255 998244353
285 998244353
121 998244353
178 998244353
120 998244353
348 998244353
114 998244353
391 998244353
237 998244353
489 998244353
150 998244353
476 998244353
13 998244353
215 998244353
430 998244353
471 998244353
286 998244353
474 998244353
330 998244353
53 998244353
97 998244353
324 998244353
372 998244353
442 998244353
153 998244353
63 998244353
382 998244353
172 998244353
460 998244353
371 998244353
500 998244353
366 998244353
258 998244353
481 998244353
497 998244353
218 998244353
261 998244353
426 998244353
468 998244353
345 998244353
99 998244353
157 998244353
`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// solve mirrors 1558B.go.
func solve(n int, mod int64) int64 {
	dp := make([]int64, n+2)
	diff := make([]int64, n+2)
	dp[1] = 1 % mod
	prefix := dp[1]
	for j := 2; j <= n; j++ {
		l := j
		r := l + j
		if r > n+1 {
			r = n + 1
		}
		diff[l] = (diff[l] + dp[1]) % mod
		diff[r] = (diff[r] - dp[1]) % mod
	}
	cur := int64(0)
	for i := 2; i <= n; i++ {
		cur = (cur + diff[i]) % mod
		if cur < 0 {
			cur += mod
		}
		dp[i] = (prefix + cur) % mod
		prefix = (prefix + dp[i]) % mod
		for j := 2; i*j <= n; j++ {
			l := i * j
			r := l + j
			if r > n+1 {
				r = n + 1
			}
			diff[l] = (diff[l] + dp[i]) % mod
			diff[r] = (diff[r] - dp[i]) % mod
		}
	}
	return dp[n] % mod
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	if len(lines)-1 != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(lines)-1)
	}
	res := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: expected 2 values", i)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i, err)
		}
		mod, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d bad mod: %v", i, err)
		}
		res = append(res, testCase{n: n, mod: mod})
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.mod)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := strconv.FormatInt(solve(tc.n, tc.mod), 10)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
