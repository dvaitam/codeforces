package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseF struct {
	n int
	a []int
	b []int
}

func parseTests(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseF, t)
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, 2*n)
		b := make([]int, 2*n)
		for j := 0; j < 2*n; j++ {
			fmt.Fscan(in, &a[j], &b[j])
		}
		cases[i] = testCaseF{n: n, a: a, b: b}
	}
	return cases, nil
}

type pair struct {
	first  int
	second bool
}

func solve(tc testCaseF) ([]string, bool) {
	n := tc.n
	m := 2 * n
	a := make([]int, m+1)
	b := make([]int, m+1)
	for i := 0; i < m; i++ {
		a[i+1] = tc.a[i]
		b[i+1] = tc.b[i]
	}
	ap := make([]int, 2*n+1)
	bp := make([]int, 2*n+1)
	vis := make([]bool, m+1)
	del := make([]bool, m+1)
	L := make([]pair, m+1)
	R := make([]pair, m+1)
	lson := make([][]int, m+1)
	rson := make([][]int, m+1)
	ans := make([]int, 0)

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
		ap[n+a[i]] = i
		bp[n+b[i]] = i
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
				return nil, false
			}
		}
	}
	if len(ans) == 0 {
		return nil, false
	}
	output := []string{"Yes"}
	for _, id := range ans {
		output = append(output, fmt.Sprintf("%d %d", a[id], b[id]))
	}
	return output, true
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expect, ok := solve(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < 2*tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
		}
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if !ok {
			if len(lines) == 0 || strings.ToUpper(strings.TrimSpace(lines[0])) != "NO" {
				fmt.Printf("case %d expected NO\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if len(lines) != len(expect) {
			fmt.Printf("case %d line count mismatch expected %d got %d\n", idx+1, len(expect), len(lines))
			os.Exit(1)
		}
		for i := range lines {
			if strings.TrimSpace(lines[i]) != expect[i] {
				fmt.Printf("case %d mismatch on line %d\nexpected: %s\n got: %s\n", idx+1, i+1, expect[i], lines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
