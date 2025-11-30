package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pair struct {
	first  int
	second bool
}

type testCase struct {
	n int
	a []int
	b []int
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
100
3
-1 3
2 -3
3 -2
-3 2
1 -1
-2 1
1
1 1
-1 -1
3
3 2
-3 -1
1 -3
-1 1
-2 3
2 -2
1
-1 -1
1 1
3
1 -3
-3 -1
-2 2
3 1
-1 -2
2 3
1
-1 -1
1 1
2
-2 -1
-1 -2
1 2
2 1
1
-1 1
1 -1
1
-1 -1
1 1
2
-2 2
2 -1
1 1
-1 -2
3
-1 2
-3 -2
1 3
-2 1
3 -1
2 -3
2
2 -2
-1 2
1 1
-2 -1
1
1 1
-1 -1
2
1 1
-1 -2
-2 2
2 -1
1
1 1
-1 -1
3
1 -1
-1 -3
2 2
-2 3
3 -2
-3 1
1
-1 1
1 -1
2
-1 -1
1 2
2 -2
-2 1
3
2 1
-3 2
1 -2
3 -3
-1 3
-2 -1
3
-3 3
3 1
-2 -1
2 -3
-1 -2
1 2
1
1 1
-1 -1
1
1 1
-1 -1
2
-1 1
1 -1
2 2
-2 -2
2
-1 -1
2 2
1 -2
-2 1
1
-1 1
1 -1
3
2 -1
-1 -2
1 1
-2 3
-3 2
3 -3
3
-1 -1
1 -2
2 1
3 2
-3 3
-2 -3
2
2 -1
1 2
-1 -2
-2 1
2
-1 2
1 -1
-2 1
2 -2
2
1 -1
-2 1
-1 2
2 -2
3
-1 1
1 3
-2 2
-3 -1
2 -2
3 -3
1
1 -1
-1 1
3
-3 1
2 -3
-1 2
1 3
-2 -1
3 -2
2
1 2
2 -1
-2 1
-1 -2
1
-1 -1
1 1
3
2 2
3 -1
1 -2
-1 1
-3 -3
-2 3
1
1 -1
-1 1
1
1 1
-1 -1
3
-1 3
3 1
1 -1
-2 -3
-3 -2
2 2
3
-1 -1
1 -2
-3 3
-2 1
2 2
3 -3
3
-2 2
1 -2
-3 -3
3 -1
-1 3
2 1
1
1 1
-1 -1
1
-1 1
1 -1
1
1 1
-1 -1
1
-1 1
1 -1
2
1 -2
-1 1
-2 2
2 -1
2
-1 -2
-2 2
2 -1
1 1
1
-1 -1
1 1
1
1 -1
-1 1
1
1 1
-1 -1
3
-3 3
1 1
2 2
3 -1
-1 -2
-2 -3
1
-1 1
1 -1
1
1 1
-1 -1
3
-1 2
2 -3
1 3
-3 -1
-2 1
3 -2
2
1 2
2 -1
-1 -2
-2 1
1
-1 1
1 -1
2
1 -2
-1 1
-2 2
2 -1
2
1 1
-2 -2
-1 -1
2 2
2
-1 2
1 1
2 -2
-2 -1
3
-2 2
-3 1
3 3
-1 -1
2 -2
1 -3
3
1 2
-3 1
-2 -1
2 -3
-1 3
3 -2
1
-1 1
1 -1
1
-1 -1
1 1
3
2 -2
-3 3
3 -1
1 -3
-1 1
-2 2
1
1 1
-1 -1
2
1 -1
-1 -2
2 2
-2 1
2
-1 2
-2 -2
1 1
2 -1
1
-1 1
1 -1
3
-2 -1
-1 3
-3 -3
1 1
3 -2
2 2
3
2 -2
1 -1
-2 3
3 2
-1 -3
-3 1
3
-1 -1
-2 3
-3 2
2 -3
3 -2
1 1
1
-1 1
1 -1
3
-1 2
-2 1
3 3
-3 -3
1 -2
2 -1
1
1 -1
-1 1
2
1 -1
-1 1
-2 -2
2 2
1
-1 1
1 -1
2
-2 2
1 -2
2 1
-1 -1
2
-2 1
-1 -1
2 2
1 -2
1
-1 -1
1 1
3
1 2
-1 3
-3 1
2 -2
-2 -3
3 -1
1
-1 1
1 -1
2
-1 1
-2 2
1 -2
2 -1
1
1 1
-1 -1
1
-1 -1
1 1
1
-1 -1
1 1
1
-1 -1
1 1
1
1 1
-1 -1
3
-1 3
-3 2
2 -3
-2 -2
3 -1
1 1
2
1 1
-1 -2
2 2
-2 -1
2
1 -1
2 -2
-1 1
-2 2
2
-1 1
1 -2
-2 2
2 -1
3
-2 -1
-1 -3
2 1
1 -2
3 3
-3 2
1
-1 1
1 -1
1
-1 -1
1 1
1
-1 -1
1 1
2
1 1
-2 -2
2 2
-1 -1
2
2 1
-1 -2
-2 -1
1 2
2
2 -2
-1 -1
-2 1
1 2
2
-1 1
2 -2
-2 -1
1 2
1
-1 -1
1 1
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		m := 2 * n
		a := make([]int, m+1)
		b := make([]int, m+1)
		for j := 1; j <= m; j++ {
			ai, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad a[%d]: %v", i+1, j, err)
			}
			bi, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad b[%d]: %v", i+1, j, err)
			}
			a[j] = ai
			b[j] = bi
		}
		res = append(res, testCase{n: n, a: a, b: b})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of data")
	}
	return res, nil
}

// solve mirrors 1503F.go and returns expected output string.
func solve(tc testCase) string {
	n := tc.n
	m := 2 * n
	a := tc.a
	b := tc.b
	ap := make([]int, 2*n+1)
	bp := make([]int, 2*n+1)
	vis := make([]bool, m+1)
	del := make([]bool, m+1)
	L := make([]pair, m+1)
	R := make([]pair, m+1)
	lson := make([][]int, m+1)
	rson := make([][]int, m+1)
	ans := make([]int, 0, m)

	for i := 1; i <= m; i++ {
		ap[n+a[i]] = i
		bp[n+b[i]] = i
	}

	var dfs func(int)
	dfs = func(id int) {
		for _, v := range lson[id] {
			dfs(v)
		}
		ans = append(ans, id)
		for _, v := range rson[id] {
			dfs(v)
		}
	}

	var calc func(int) bool
	calc = func(start int) bool {
		queue := make([]int, 0)
		cur := start
		for {
			next := R[cur].first
			if R[cur].second == R[next].second {
				queue = append(queue, cur)
			}
			cur = next
			if cur == start {
				break
			}
		}
		for L[start].first != R[start].first || L[start].second != R[start].second {
			if R[start].first == L[start].first || len(queue) == 0 {
				return false
			}
			i := queue[0]
			queue = queue[1:]
			if del[i] {
				continue
			}
			ni := R[i].first
			if R[i].second != R[ni].second {
				continue
			}
			del[ni] = true
			ni2 := R[ni].first
			del[ni2] = true
			if R[i].second {
				rson[i] = append(rson[i], ni)
				rson[i] = append(rson[i], ni2)
			} else {
				lson[i] = append([]int{ni}, lson[i]...)
				lson[i] = append([]int{ni2}, lson[i]...)
			}
			j := R[ni2].first
			L[j].first = i
			R[i].first = j
			R[i].second = !L[j].second
			start = i
			if L[i].second != R[i].second {
				queue = append(queue, L[i].first)
			}
			if R[i].second == R[R[i].first].second {
				queue = append(queue, i)
			}
		}
		if R[start].second {
			dfs(start)
			dfs(R[start].first)
		} else {
			dfs(R[start].first)
			dfs(start)
		}
		return true
	}

	for i := 1; i <= m; i++ {
		if !vis[i] {
			j := i
			for !vis[j] {
				vis[j] = true
				nxa := ap[n-a[j]]
				nxb := bp[n-b[j]]
				if !vis[nxa] {
					R[j] = pair{nxa, a[j] > 0}
					L[nxa] = pair{j, a[j] < 0}
					j = nxa
				} else {
					R[j] = pair{nxb, b[j] > 0}
					L[nxb] = pair{j, b[j] < 0}
					j = nxb
				}
			}
			if !calc(i) {
				return "NO"
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	for _, id := range ans {
		fmt.Fprintf(&sb, "%d %d\n", a[id], b[id])
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 1; i <= 2*tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d: %v", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
