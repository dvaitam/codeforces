package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1322DSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

// Greedy attempt for problem described in problemD.txt.
// We process candidates in order and recruit a candidate only
// if their addition immediately increases total profit while
// respecting the non-increasing aggressiveness constraint.

func simulateGain(cnt []int, level int, c []int) int {
	gain := c[level]
	cnt[level]++
	for lvl := level; cnt[lvl] == 2; lvl++ {
		cnt[lvl] = 0
		if lvl+1 >= len(cnt) {
			cnt = append(cnt, 0)
		}
		cnt[lvl+1]++
		gain += c[lvl+1]
	}
	return gain
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	l := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &l[i])
	}
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	c := make([]int, n+m+5)
	for i := 1; i <= n+m; i++ {
		fmt.Fscan(in, &c[i])
	}

	cnt := make([]int, n+m+5)
	profit := 0
	cost := 0
	maxAllowed := n + m + 5

	for i := 0; i < n; i++ {
		if l[i] > maxAllowed {
			continue
		}
		tmp := make([]int, len(cnt))
		copy(tmp, cnt)
		g := simulateGain(tmp, l[i], c)
		if profit+g-(cost+s[i]) > profit-cost {
			cnt = tmp
			profit += g
			cost += s[i]
			if l[i] < maxAllowed {
				maxAllowed = l[i]
			}
		}
	}

	fmt.Fprintln(out, profit-cost)
}
`

// Keep the embedded reference solution reachable.
var _ = solution1322DSource

type testCase struct {
	n int
	m int
	l []int
	s []int
	c []int
}

const testcasesRaw = `4 4 2 4 2 2 4 0 2 1 -1 -4 2 2 -2 5 2 -1
4 2 2 1 1 2 2 3 2 3 -3 -3 4 3 2 -4
2 5 3 5 2 0 -1 2 -2 1 -2 -1 4
3 1 1 1 1 5 0 4 5 3 -5 -5
1 4 3 3 0 5 -5 3 1
1 4 3 5 -5 0 -5 -5 0
1 5 1 5 3 5 -4 4 3 -5
2 5 5 4 2 0 -4 -3 -4 -3 3 -2 -2
4 5 3 4 2 1 4 4 3 0 -2 -3 1 5 2 -4 -2 -2 -4
2 4 4 4 2 4 -2 3 -5 3 5 1
2 4 2 1 2 3 4 2 5 2 4 1
3 5 2 1 3 4 4 1 0 -1 -2 3 1 -1 -1 1
3 3 1 1 2 2 1 0 -2 1 3 -3 0 -2
4 4 2 3 3 4 0 4 3 3 5 -4 -3 5 -4 3 -1 -1
1 4 4 1 -5 4 -3 -5 1
4 3 2 3 2 1 2 3 5 4 4 -5 -4 -5 -3 4 3
2 2 2 2 1 4 -1 2 -2 5
1 3 1 4 -2 -2 -3 -5
1 4 3 2 3 -5 1 2 4
2 1 1 1 4 3 0 -4 0
1 4 2 2 1 -5 1 2 2
2 2 1 1 4 1 5 3 1 0
3 2 1 2 2 0 5 4 4 -3 -4 2 -5
4 3 2 2 1 1 5 4 1 1 1 -2 -1 4 3 4 -5
1 2 1 4 5 0 -1
4 5 3 3 5 3 2 1 0 5 2 3 -4 -4 -3 -3 0 1 -2
5 3 3 2 3 2 2 4 5 0 4 0 5 1 5 -4 4 1 2 -4
1 4 2 3 -3 1 3 -5 0
4 3 2 1 1 3 4 3 4 2 4 -2 3 5 5 5 -5
2 4 4 2 4 4 -5 -2 -1 3 -2 -1
1 5 3 5 -1 1 -3 -5 -2 4
3 2 2 2 1 0 3 3 3 -5 -3 5 -4
3 4 2 3 4 4 2 4 -3 -4 0 -2 5 -5 1
4 3 2 1 3 3 1 4 4 4 -5 2 -2 -2 -1 -3 -3
4 2 1 1 1 2 3 3 5 0 -3 4 -1 -1 5 -5
5 2 2 2 2 1 2 0 5 3 2 1 4 4 -5 -4 0 5 2
2 3 2 3 0 0 -4 5 4 1 -2
4 2 2 2 2 1 1 0 4 3 -4 1 5 0 -3 -4
5 5 3 4 5 2 4 5 3 0 1 2 -1 3 -1 5 1 4 -5 5 2 -4
4 3 3 2 3 3 2 3 4 0 5 -1 -2 1 2 0 2
2 4 3 1 2 5 4 -2 4 -3 5 3
3 1 1 1 1 4 3 2 -3 -5 5 1
3 3 2 1 3 0 0 1 1 5 -4 -5 1 2
5 2 1 2 1 1 2 4 0 0 0 2 -3 2 -3 -2 -2 5 -2
5 1 1 1 1 1 1 5 0 0 5 4 -5 -5 5 -2 1 -1
1 4 3 0 3 -3 -4 0 0
5 2 1 2 1 1 1 2 0 0 2 3 1 3 -1 3 -3 -4 0
4 1 1 1 1 1 4 2 5 5 4 0 -2 -5 -5
1 1 1 2 -4 5
1 3 2 3 -5 -4 1 -2
3 3 3 3 3 5 5 0 -3 3 -4 0 0 -1
3 1 1 1 1 2 2 0 -3 2 0 -2
2 1 1 1 4 5 -4 -3 -3
5 2 1 2 2 1 2 1 3 4 4 1 2 3 -4 2 -2 0 4
3 3 3 3 3 4 0 3 5 -4 0 1 1 0
2 2 2 2 3 5 3 -1 5 -2
2 1 1 1 1 1 -5 0 1
1 1 1 2 -3 -3
3 4 4 3 3 2 2 3 -3 2 5 -3 5 2 -3
4 5 3 3 1 1 5 1 5 2 0 5 -2 -4 -1 0 4 -5 -4
4 5 5 2 1 4 3 0 1 0 -2 -2 5 4 -2 0 0 -5 -1
3 1 1 1 1 1 3 0 5 1 3 -3
4 4 3 1 3 1 0 3 3 5 -5 -2 3 -5 5 -4 1 -4
4 2 2 2 2 1 4 0 4 1 -3 4 -3 2 -3 2
2 1 1 1 3 3 -1 5 2
3 5 5 5 1 1 4 5 5 5 -5 2 2 2 -1 1
2 1 1 1 3 3 -2 4 4
1 5 2 4 0 -2 -3 -1 3 3
2 4 1 3 3 1 -4 4 0 1 -3 5
4 5 4 1 3 4 2 2 3 2 1 5 -1 -5 -2 2 0 -2 -4
3 1 1 1 1 4 3 1 5 5 5 5
4 3 1 1 1 3 2 0 1 2 -4 -2 -3 -2 4 1 5
3 3 1 1 3 4 3 0 -3 2 5 4 -4 -5
3 4 3 2 2 5 5 2 -1 -3 -3 0 1 -1 0
4 4 4 2 4 2 2 3 5 3 2 -5 0 2 0 5 -5 0
4 3 3 3 2 1 0 5 4 2 5 -1 4 1 0 5 5
4 1 1 1 1 1 3 4 2 5 -5 0 -5 4 2
5 5 4 2 1 3 5 4 1 4 4 1 5 -4 3 1 -4 3 -5 -1 0 0
1 4 2 4 4 -5 4 5 4
5 4 3 4 3 2 3 1 1 1 3 5 -2 -2 -1 0 -2 4 0 -3 -4
4 2 1 1 2 1 3 5 2 4 -4 -4 -1 -2 -5 -5
3 1 1 1 1 3 3 2 -2 5 3 0
5 5 1 3 4 3 1 4 1 4 3 4 4 1 3 0 -1 3 -3 5 -3 -4
2 1 1 1 5 5 -2 5 3
3 4 3 1 4 0 1 4 1 -5 1 5 1 4 3
5 2 1 1 2 2 2 5 1 3 3 1 5 5 5 4 4 -1 3
5 4 3 1 1 3 3 3 1 4 2 0 -4 2 -3 2 5 -1 5 3 -3
2 5 1 2 4 5 -3 2 -5 -4 3 0 0
4 3 1 1 2 3 3 5 4 4 3 -4 4 4 4 2 -1
1 3 3 3 -4 5 5 3
3 1 1 1 1 2 5 1 0 2 -2 1
3 5 1 3 3 4 3 1 2 1 4 4 0 1 1 2
2 2 2 1 1 0 2 -5 -3 -4
2 5 2 1 5 0 1 -2 3 -5 -3 1 -3
1 4 3 4 -2 5 2 -5 -4
1 2 1 0 3 -4 4
1 1 1 2 2 4
5 2 2 2 1 2 2 1 5 5 3 1 0 5 -2 -2 -3 4 1
4 4 1 2 1 3 0 0 3 0 -3 -2 -1 -4 5 -2 -5 4
4 4 4 4 2 3 5 5 0 5 0 2 1 0 4 -5 5 -5`

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var res []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		pos := 2
		if len(fields) < 2+n+n+(n+m) {
			continue
		}
		l := make([]int, n)
		for i := 0; i < n; i++ {
			l[i], _ = strconv.Atoi(fields[pos+i])
		}
		pos += n
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i], _ = strconv.Atoi(fields[pos+i])
		}
		pos += n
		c := make([]int, n+m+1)
		for i := 1; i <= n+m; i++ {
			c[i], _ = strconv.Atoi(fields[pos+i-1])
		}
		res = append(res, testCase{n: n, m: m, l: l, s: s, c: c})
	}
	return res
}

func simulateGain(cnt []int, level int, c []int) int {
	gain := c[level]
	cnt[level]++
	for lvl := level; cnt[lvl] == 2; lvl++ {
		cnt[lvl] = 0
		if lvl+1 >= len(cnt) {
			cnt = append(cnt, 0)
		}
		cnt[lvl+1]++
		gain += c[lvl+1]
	}
	return gain
}

func solveExpected(tc testCase) int {
	cnt := make([]int, tc.n+tc.m+5)
	profit := 0
	cost := 0
	maxAllowed := tc.n + tc.m + 5
	for i := 0; i < tc.n; i++ {
		if tc.l[i] > maxAllowed {
			continue
		}
		tmp := make([]int, len(cnt))
		copy(tmp, cnt)
		g := simulateGain(tmp, tc.l[i], tc.c)
		if profit+g-(cost+tc.s[i]) > profit-cost {
			cnt = tmp
			profit += g
			cost += tc.s[i]
			if tc.l[i] < maxAllowed {
				maxAllowed = tc.l[i]
			}
		}
	}
	return profit - cost
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i := 1; i <= tc.n+tc.m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", tc.c[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	expect := solveExpected(tc)
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	gotStr := strings.TrimSpace(string(out))
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("case %d failed: invalid output %q\ninput:\n%s", idx, gotStr, input)
	}
	if got != expect {
		return fmt.Errorf("case %d failed: expected %d got %d\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
